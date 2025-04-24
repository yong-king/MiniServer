package sequence

type Sequence interface{
	Next() (uint64, error)
}

type LinkMapping interface {
    SetShortLink(shortLink, longLink string) error
    GetLongLink(shortLink string) (string, error)
}
