package elasticsearch

func getIndexMapping(indexName string) string {
	mapping := ``
	if indexName == "url-metadata" {
		mapping = `
		{
			"mappings": {
				"properties": {
					"url": {
						"type": "text",
						"fields": {
							"keyword": {
								"type": "keyword",
								"ignore_above": 256
							}
						}
					},
					"shortened_url": {
						"type": "text",
						"fields": {
							"keyword": {
								"type": "keyword",
								"ignore_above": 256
							}
						}
					},
					"expiry": {
						"type": "date"
					}
				}
			}
		}
		`
	}
	return mapping
}
