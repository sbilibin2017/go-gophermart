package queries

var OrderSaveQuery = `
	INSERT INTO orders (number, status, accrual, updated_at)
	VALUES (:number, :status, :accrual, now())
	ON CONFLICT (number) DO UPDATE 
	SET 
		status = EXCLUDED.status,
		accrual = EXCLUDED.accrual,
		updated_at = now()
`
