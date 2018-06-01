package userConfig

var disHonestList = map[int64]string{1434: ""}

var notInterestedList = map[int64]string{}

func FilterCompany(id int64) bool {
	if _, ok := disHonestList[id]; ok {
		return false
	}
	if _, ok := notInterestedList[id]; ok {
		return false
	}

	return true
}
