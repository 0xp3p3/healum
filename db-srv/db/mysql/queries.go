package mysql

var (
	selectQ      = "SELECT id, created, updated, name, parameter1, parameter2, parameter3, lat, lng, metadata, lng from %s.%s"
	mysqlQueries = map[string]string{
		"delete": "DELETE from %s.%s where id = ? and parameter3 = ? limit 1",
		"create": `INSERT into %s.%s (id, created, updated, name, parameter1, parameter2, parameter3, lat, lng, metadata) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"update": "UPDATE %s.%s set updated = ?, name = ?, parameter1 = ?, parameter2 = ?, parameter3 = ?, lat = ?, lng = ?, metadata = ? where id = ?",
		"read":   selectQ + " where id = ? and parameter3 = ? limit 1",
		"readId":   selectQ + " where id = ? limit 1",

		"searchAsc":  selectQ + " where created >= ? and created <= ? order by created asc limit ? offset ?",
		"searchDesc": selectQ + " where created >= ? and created <= ? order by created desc limit ? offset ?",

		"searchParameter3Asc":  selectQ + " where created >= ? and created <= ? and parameter3 = ? order by created asc limit ? offset ?",
		"searchParameter3Desc": selectQ + " where created >= ? and created <= ? and parameter3 = ? order by created desc limit ? offset ?",
		"sNameAsc":   selectQ + " where created >= ? and created <= ? and name = ? order by created asc limit ? offset ?",
		"sNameDesc":  selectQ + " where created >= ? and created <= ? and name = ? order by created desc limit ? offset ?",

		"sParameter1Asc":  selectQ + " where created >= ? and created <= ? and parameter1 = ? order by created asc limit ? offset ?",
		"sParameter1Desc": selectQ + " where created >= ? and created <= ? and parameter1 = ? order by created desc limit ? offset ?",

		"sParameter2Asc":  selectQ + " where created >= ? and created <= ? and parameter2 like ? order by created asc limit ? offset ?",
		"sParameter2Desc": selectQ + " where created >= ? and created <= ? and parameter2 like ? order by created desc limit ? offset ?",


		"sNameParameter1Asc":   selectQ + " where created >= ? and created <= ? and name = ? and parameter1 = ? order by created asc limit ? offset ?",
		"sNameParameter1Desc":  selectQ + " where created >= ? and created <= ? and name = ? and parameter1 = ? order by created desc limit ? offset ?",

		"sNameParameter2Asc":   selectQ + " where created >= ? and created <= ? and name = ? and parameter2 = ? order by created asc limit ? offset ?",
		"sNameParameter2Desc":  selectQ + " where created >= ? and created <= ? and name = ? and parameter2 = ? order by created desc limit ? offset ?",

		"sNameParameter3Asc":   selectQ + " where created >= ? and created <= ? and parameter3 = ? and name = ? order by created asc limit ? offset ?",
		"sNameParameter3Desc":  selectQ + " where created >= ? and created <= ? and parameter3 = ? and name = ? order by created desc limit ? offset ?",

		"sParameter1Parameter3Asc":  selectQ + " where created >= ? and created <= ? and parameter3 = ? and parameter1 = ? order by created asc limit ? offset ?",
		"sParameter1Parameter3Desc": selectQ + " where created >= ? and created <= ? and parameter3 = ? and parameter1 = ? order by created desc limit ? offset ?",

		"sParameter1Parameter2Asc":  selectQ + " where created >= ? and created <= ? and parameter2 = ? and parameter1 = ? order by created asc limit ? offset ?",
		"sParameter1Parameter2Desc": selectQ + " where created >= ? and created <= ? and parameter2 = ? and parameter1 = ? order by created desc limit ? offset ?",

		"sParameter2Parameter3Asc":  selectQ + " where created >= ? and created <= ? and parameter3 = ? and parameter2 like ? order by created asc limit ? offset ?",
		"sParameter2Parameter3Desc": selectQ + " where created >= ? and created <= ? and parameter3 = ? and parameter2 like ? order by created desc limit ? offset ?",

		"sNameAndParameter1Parameter3Asc":  selectQ + " where created >= ? and created <= ? and parameter3 = ? and name = ? and parameter1 = ? order by created asc limit ? offset ?",
		"sNameAndParameter1Parameter3Desc": selectQ + " where created >= ? and created <= ? and parameter3 = ? and name = ? and parameter1 = ? order by created desc limit ? offset ?",

		"sNameAndParameter1Asc":  selectQ + " where created >= ? and created <= ? and name = ? and parameter1 = ? order by created asc limit ? offset ?",
		"sNameAndParameter1Desc": selectQ + " where created >= ? and created <= ? and name = ? and parameter1 = ? order by created desc limit ? offset ?",

		"sNameAndParameter2Asc":  selectQ + " where created >= ? and created <= ? and name = ? and parameter2 = ? order by created asc limit ? offset ?",
		"sNameAndParameter2Desc": selectQ + " where created >= ? and created <= ? and name = ? and parameter2 = ? order by created desc limit ? offset ?",

		"sNameAndParameter1Parameter2Asc":  selectQ + " where created >= ? and created <= ? and parameter2 = ? and name = ? and parameter1 = ? order by created asc limit ? offset ?",
		"sNameAndParameter1Parameter2Desc": selectQ + " where created >= ? and created <= ? and parameter2 = ? and name = ? and parameter1 = ? order by created desc limit ? offset ?",

		"sNameAndParameter1Parameter2Parameter3Asc":  selectQ + " where created >= ? and created <= ? and parameter3 = ? and name = ? and parameter1 = ? and parameter2 = ? order by created asc limit ? offset ?",
		"sNameAndParameter1Parameter2Parameter3Desc": selectQ + " where created >= ? and created <= ? and parameter3 = ? and name = ? and parameter1 = ? and parameter2 = ? order by created desc limit ? offset ?",

		"distanceSearch": "SELECT id, created, updated, name, parameter1, parameter2, parameter3, lat, lng, metadata, (3959 * acos(cos(radians(?)) * cos(radians(lat)) * cos( radians(lng) - radians(?)) + sin(radians(?)) * sin(radians(lat)))) AS distance FROM %s.%s HAVING distance <= ? ORDER BY distance",
	}

	searchMetadataQ = selectQ + " where metadata like ?"
)
