package gochi

// Unimplemented

// func TestStorage(t *testing.T) {
// 	g := New()
//
// 	text := "test"
//
// 	writeData := Storage{
// 		Path: "test/test.txt",
// 		Body: []byte(text),
// 	}
//
// 	getData := Storage{
// 		Path: "test.txt",
// 	}
//
// 	_, ctx, spindown := SpinUp(t)
// 	defer spindown()
//
// 	err := writeData.Write(ctx)
// 	g.Ok(t, err)
//
// 	err = getData.Get(ctx)
// 	g.Ok(t, err)
// 	g.Equals(t, text, string(getData.Body))
// }
