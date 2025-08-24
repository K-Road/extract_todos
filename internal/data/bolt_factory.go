package data

// func BoltFactory(dbFile string, readOnly bool) func() (config.DataProvider, error) {
// 	return func() (config.DataProvider, error) {
// 		bp := &BoltProvider{}
// 		var opts *bolt.Options
// 		if readOnly {
// 			opts = &bolt.Options{ReadOnly: true}
// 		} else {
// 			opts = nil
// 		}
// 		if err := bp.OpenDBWithOptions(dbFile, opts); err != nil {
// 			return nil, err
// 		}
// 		return bp, nil
// 	}
// }
