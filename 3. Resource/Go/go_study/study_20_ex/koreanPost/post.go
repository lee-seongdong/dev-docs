package koreanPost

import "fmt"

type PostSender struct {
	// ...
}

func (k *PostSender) Send(parcel string) {
	fmt.Printf("우체국 택배 %s 발송\n", parcel)
}
