package queries

var RewardFilterQuery = `
	SELECT match, reward, reward_type
	FROM rewards
	WHERE match ILIKE ANY($1)
`
