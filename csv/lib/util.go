package wcsv

func deepCopy(l1 []ColumnDef) []ColumnDef {
	var l2 []ColumnDef
	var count = len(l2)
	for i := 0; i < count; i++ {
		l2[i].Name = l1[i].Name
		l2[i].Required = l1[i].Required
		l2[i].CaseSensitive = l1[i].CaseSensitive
		l2[i].CanonicalIndex = l1[i].CanonicalIndex
		l2[i].Index = l1[i].Index
	}
	return l2
}
