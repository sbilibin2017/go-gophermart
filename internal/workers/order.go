package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrderListRepository interface {
	ListOrdered(ctx context.Context, login *string) ([]*types.UserOrderDB, error)
}

type UserOrderSaveRepository interface {
	Save(ctx context.Context, number string, login string, status string, accrual *int64) error
}

type UserBalanceWithdrawSaveRepository interface {
	Save(ctx context.Context, login string, number string, sum int64) error
}

type OrderResult struct {
	Order *types.UserOrderDB
	Err   error
}

func StartOrderWorkerDaemon(
	ctx context.Context,
	uol UserOrderListRepository,
	uos UserOrderSaveRepository,
	ubws UserBalanceWithdrawSaveRepository,
	client *resty.Client,
	accrualSystemAddress string,
	maxConcurrentRequests int,
	numWorkers int,
	tickerInterval time.Duration,
) error {
	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				select {
				case <-doneCh:
					doneCh = make(chan struct{})
					if err := RunOrderWorkers(ctx, uol, uos, ubws, client, accrualSystemAddress, maxConcurrentRequests, numWorkers, doneCh); err != nil {
						log.Printf("Ошибка при запуске воркеров: %v", err)
					}
				default:
					log.Println("Пропуск запуска воркеров, так как предыдущие еще не завершены.")
				}
			}
		}
	}()

	<-ctx.Done()
	return nil
}

func RunOrderWorkers(
	ctx context.Context,
	uol UserOrderListRepository,
	uos UserOrderSaveRepository,
	ubws UserBalanceWithdrawSaveRepository,
	client *resty.Client,
	accrualSystemAddress string,
	semaPool int,
	numWorkers int,
	doneCh chan struct{},
) error {
	sema := newSemaphore(semaPool)

	orderCh := orderGenerator(ctx, uol)

	updatedOrderCh := make(chan OrderResult)

	// Fix: remove ctx from the arguments, it should not be passed to updateOrderStageWorkerPool
	go updateOrderStageWorkerPool(accrualSystemAddress, client, orderCh, updatedOrderCh, sema, numWorkers)

	savedOrderCh := make(chan OrderResult)
	go saveOrderStageWorkerPool(ctx, uos, updatedOrderCh, savedOrderCh, sema, numWorkers)

	finalBalanceCh := make(chan OrderResult)
	go saveBalanceStageWorkerPool(ctx, ubws, savedOrderCh, finalBalanceCh, sema, numWorkers)

	fanInCh := fanIn(doneCh, sema, finalBalanceCh)

	for res := range fanInCh {
		if res.Err != nil {
			log.Printf("Ошибка обработки заказа %s: %v", res.Order.Number, res.Err)
		} else {
			log.Printf("Успешно обработан заказ %s", res.Order.Number)
		}
	}

	close(doneCh)

	return nil
}

func fanIn(doneCh <-chan struct{}, sema *semaphore, resultChs ...chan OrderResult) chan OrderResult {
	finalCh := make(chan OrderResult)
	var wg sync.WaitGroup

	for _, ch := range resultChs {
		chClosure := ch
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer func() {
				sema.Release()
			}()

			for res := range chClosure {
				select {
				case <-doneCh:
					return
				case finalCh <- res:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(finalCh)
	}()

	return finalCh
}

type semaphore struct {
	semaCh chan struct{}
}

func newSemaphore(maxReq int) *semaphore {
	return &semaphore{
		semaCh: make(chan struct{}, maxReq),
	}
}

func (s *semaphore) Acquire() {
	s.semaCh <- struct{}{}
}

func (s *semaphore) Release() {
	<-s.semaCh
}

func orderGenerator(
	ctx context.Context,
	uol UserOrderListRepository,
) chan *types.UserOrderDB {
	out := make(chan *types.UserOrderDB)

	go func() {
		defer close(out)
		orders, err := uol.ListOrdered(ctx, nil)
		if err != nil {
			err = fmt.Errorf("ошибка получения заказов: %v", err)
			log.Println(err)
			return
		}

		for _, order := range orders {
			select {
			case <-ctx.Done():
				log.Println("orderGenerator: остановка по сигналу контекста")
				return
			case out <- order:
			}
		}
	}()

	return out
}

func updateOrderStage(
	accrualSystemAddress string,
	client *resty.Client,
	in <-chan *types.UserOrderDB,
	out chan<- OrderResult,
	sema *semaphore,
) {
	defer close(out)

	for order := range in {
		sema.Acquire()

		func() {
			defer sema.Release()
			url := fmt.Sprintf("%s/api/orders/%s", accrualSystemAddress, order.Number)
			resp, err := client.R().Get(url)

			if err != nil {
				log.Printf("Ошибка при запросе к сервису начислений для заказа %s: %v", order.Number, err)
				out <- OrderResult{Order: order, Err: err}
				return
			}

			if resp.StatusCode() != http.StatusOK {
				err = fmt.Errorf("ошибка запроса: статус код %d для заказа %s", resp.StatusCode(), order.Number)
				log.Println(err)
				out <- OrderResult{Order: order, Err: err}
				return
			}

			var orderAccrualResponse types.OrderGetResponse

			if err := json.Unmarshal(resp.Body(), &orderAccrualResponse); err != nil {
				err = fmt.Errorf("ошибка десериализации ответа для заказа %s: %v", order.Number, err)
				log.Println(err)
				out <- OrderResult{Order: order, Err: err}
				return
			}

			switch orderAccrualResponse.Status {
			case types.OrderAccrualStatusProcessed:
				order.Status = types.GophermartUserOrderStatusProcessed
				order.Accrual = &orderAccrualResponse.Accrual
			case types.OrderAccrualStatusProcessing:
				order.Status = types.GophermartUserOrderStatusProcessing
			case types.OrderAccrualStatusInvalid:
				order.Status = types.GophermartUserOrderStatusInvalid
			default:
				err = fmt.Errorf("неизвестный статус заказа %s: %s", order.Number, orderAccrualResponse.Status)
				log.Println(err)
				out <- OrderResult{Order: order, Err: err}
				return
			}

			out <- OrderResult{Order: order, Err: nil}
		}()
	}
}

func updateOrderStageWorkerPool(
	accrualSystemAddress string,
	client *resty.Client,
	in <-chan *types.UserOrderDB,
	out chan<- OrderResult,
	sema *semaphore,
	numWorkers int,
) {
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer sema.Release()
			updateOrderStage(accrualSystemAddress, client, in, out, sema)
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
}

func saveOrderStage(
	ctx context.Context,
	uos UserOrderSaveRepository,
	in <-chan OrderResult,
	out chan<- OrderResult,
	sema *semaphore,
) {
	defer close(out)

	for res := range in {
		sema.Acquire()

		func() {
			defer sema.Release()

			if res.Err != nil {
				log.Printf("saveOrderStage: пропуск заказа %s из-за ошибки на предыдущем этапе: %v", res.Order.Number, res.Err)
				return
			}

			err := uos.Save(ctx, res.Order.Number, res.Order.Login, res.Order.Status, res.Order.Accrual)
			if err != nil {
				err = fmt.Errorf("saveOrderStage: ошибка сохранения заказа %s: %w", res.Order.Number, err)
				log.Println(err)
				out <- OrderResult{Order: res.Order, Err: err}
				return
			}

			out <- OrderResult{Order: res.Order, Err: nil}
		}()
	}
}

func saveOrderStageWorkerPool(
	ctx context.Context,
	uos UserOrderSaveRepository,
	in <-chan OrderResult,
	out chan<- OrderResult,
	sema *semaphore,
	numWorkers int,
) {
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer sema.Release()

			saveOrderStage(ctx, uos, in, out, sema)
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
}

func saveBalanceStage(
	ctx context.Context,
	ubws UserBalanceWithdrawSaveRepository,
	in <-chan OrderResult,
	sema *semaphore,
) {
	for res := range in {
		sema.Acquire()
		func() {
			defer sema.Release()

			if res.Err != nil {
				log.Printf("saveBalanceStage: ошибка на предыдущем этапе, заказ %s пропущен: %v", res.Order.Number, res.Err)
				return
			}

			if res.Order.Status == types.GophermartUserOrderStatusProcessed && res.Order.Accrual != nil {
				err := ubws.Save(ctx, res.Order.Login, res.Order.Number, *res.Order.Accrual)
				if err != nil {
					err = fmt.Errorf("saveBalanceStage: ошибка сохранения начисления по заказу %s: %w", res.Order.Number, err)
					log.Println(err)
				}
			}
		}()
	}
}

func saveBalanceStageWorkerPool(
	ctx context.Context,
	ubws UserBalanceWithdrawSaveRepository,
	in <-chan OrderResult,
	out chan<- OrderResult,
	sema *semaphore,
	numWorkers int,
) {
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer sema.Release()
			saveBalanceStage(ctx, ubws, in, sema)
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
}
