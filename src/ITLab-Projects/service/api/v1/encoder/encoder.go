package encoder

import (
	"github.com/ITLab-Projects/service/responce"
	"net/http"
	"context"
)

func EncodeResponce(
	ctx context.Context, 
	w http.ResponseWriter, 
	resp interface{},
) error {
	httpresp := resp.(responce.HTTPResponce)
	httpresp.Headers(ctx, w)
	w.WriteHeader(httpresp.StatusCode())
	return httpresp.Encode(w)
}