package service_interface

import (
	"context"

	dto_content "github.com/KaueTTS/streaming_api/src/api/v1/dto/content"
)

type ContentServiceInterface interface {
	ListContents(ctx context.Context, request dto_content.ContentListRequestDto) (dto_content.ContentResponseDto, error)
	SearchContents(ctx context.Context, request dto_content.ContentSearchRequestDto) (dto_content.ContentResponseDto, error)
}
