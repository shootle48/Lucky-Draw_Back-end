package prize

import (
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Create(ctx *gin.Context) {
	var body request.CreatePrize

	if err := ctx.Bind(&body); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	data, mserr, err := ctl.Service.Create(ctx, body)
	if err != nil {
		ms := "internal server error"
		if mserr {
			ms = err.Error()
		}
		logger.Errf(err.Error())
		response.InternalError(ctx, ms)
		return
	}

	response.Success(ctx, data)
}

func (ctl *Controller) Update(ctx *gin.Context) {
	ID := request.GetByIDPrize{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	body := request.UpdatePrize{}
	if err := ctx.Bind(&body); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	_, mserr, err := ctl.Service.Update(ctx, body, ID)
	if err != nil {
		ms := "internal server error"
		if mserr {
			ms = err.Error()
		}
		logger.Errf(err.Error())
		response.InternalError(ctx, ms)
		return
	}

	response.Success(ctx, nil)
}

func (ctl *Controller) List(ctx *gin.Context) {
	req := request.ListPrize{}
	if err := ctx.Bind(&req); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.Page == 0 {
		req.Page = 10
	}

	if req.OrderBy == "" {
		req.OrderBy = "asc"
	}

	if req.SortBy == "" {
		req.SortBy = "created_at"
	}

	data, total, err := ctl.Service.List(ctx, req)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalError(ctx, err.Error())
		return
	}
	response.SuccessWithPaginate(ctx, data, req.Size, req.Page, total)
}

func (ctl *Controller) Get(ctx *gin.Context) {
	ID := request.GetByIDPrize{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	data, err := ctl.Service.Get(ctx, ID)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalError(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

func (ctl *Controller) Delete(ctx *gin.Context) {
	ID := request.GetByIDPrize{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	err := ctl.Service.Delete(ctx, ID)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalError(ctx, err.Error())
		return
	}
	response.Success(ctx, nil)
}

// new function
func (ctl *Controller) UploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		logger.Errf("No file uploaded: %v", err)
		response.BadRequest(ctx, "กรุณาเลือกไฟล์รูปภาพ")
		return
	}

	// เปิดไฟล์
	src, err := file.Open()
	if err != nil {
		logger.Errf("Cannot open uploaded file: %v", err)
		response.InternalError(ctx, "ไม่สามารถเปิดไฟล์ได้")
		return
	}
	defer src.Close()

	// สร้าง Cloudinary client
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		logger.Errf("Cloudinary config error: %v", err)
		response.InternalError(ctx, "การตั้งค่า Cloudinary ไม่ถูกต้อง")
		return
	}

	// กำหนด timeout สำหรับ upload
	uploadCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// อัปโหลดไปยัง Cloudinary
	uploadResult, err := cld.Upload.Upload(uploadCtx, src, uploader.UploadParams{
		Folder:   "prizes",
		PublicID: fmt.Sprintf("prize_%d", time.Now().UnixNano()), // ตั้งชื่อให้ unique
	})
	if err != nil {
		logger.Errf("Upload to Cloudinary failed: %v", err)
		response.InternalError(ctx, "ไม่สามารถอัปโหลดรูปภาพได้")
		return
	}

	// ส่งกลับ URL
	response.Success(ctx, gin.H{
		"url": uploadResult.SecureURL,
	})
}
