package queries

// Обновлённый запрос, учитывающий правильные имена столбцов
var OrderExistsQuery = `SELECT EXISTS(SELECT 1 FROM orders WHERE number = $1)`
