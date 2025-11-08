package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ListProducts 获取商品列表
func ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	_ = c.Query("keyword")

	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	// 实际项目中应该查询数据库
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list": []gin.H{
				{
					"id":    1,
					"name":  "iPhone 14 Pro",
					"price": 7999.00,
					"image": "https://example.com/iphone14pro.jpg",
					"stock": 100,
					"sales": 200,
				},
				{
					"id":    2,
					"name":  "Samsung Galaxy S23",
					"price": 6999.00,
					"image": "https://example.com/s23.jpg",
					"stock": 80,
					"sales": 150,
				},
			},
			"pagination": gin.H{
				"page":  page,
				"size":  size,
				"total": 100,
			},
		},
	})
}

// GetProductDetail 获取商品详情
func GetProductDetail(c *gin.Context) {
	id := c.Param("id")

	// 实际项目中应该查询数据库
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id":            id,
			"name":          "iPhone 14 Pro",
			"price":         7999.00,
			"originalPrice": 8999.00,
			"images": []string{
				"https://example.com/iphone14pro-1.jpg",
				"https://example.com/iphone14pro-2.jpg",
			},
			"stock":        100,
			"sales":        200,
			"categoryId":   1,
			"categoryName": "手机",
			"description":  "iPhone 14 Pro 搭载A16芯片，全新灵动岛设计...",
			"specs": []gin.H{
				{
					"name":    "颜色",
					"options": []string{"黑色", "白色", "金色", "紫色"},
				},
				{
					"name":    "存储",
					"options": []string{"128GB", "256GB", "512GB", "1TB"},
				},
			},
		},
	})
}

// ListCategories 获取分类列表
func ListCategories(c *gin.Context) {
	// 实际项目中应该查询数据库
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": []gin.H{
			{
				"id":       1,
				"name":     "手机",
				"icon":     "https://example.com/icon-phone.png",
				"subCount": 10,
			},
			{
				"id":       2,
				"name":     "电脑",
				"icon":     "https://example.com/icon-laptop.png",
				"subCount": 8,
			},
		},
	})
}

// ListProductsByCategory 按分类获取商品
func ListProductsByCategory(c *gin.Context) {
	categoryId := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	// 实际项目中应该查询数据库
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"categoryId":   categoryId,
			"categoryName": "手机",
			"list": []gin.H{
				{
					"id":    1,
					"name":  "iPhone 14 Pro",
					"price": 7999.00,
					"image": "https://example.com/iphone14pro.jpg",
					"stock": 100,
					"sales": 200,
				},
				// 更多商品...
			},
			"pagination": gin.H{
				"page":  page,
				"size":  size,
				"total": 50,
			},
		},
	})
}

// CreateProduct 创建商品（管理员）
func CreateProduct(c *gin.Context) {
	// 实际项目中应该验证并创建商品
	name := c.PostForm("name")
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "商品创建成功",
		"data": gin.H{
			"id":    10,
			"name":  name,
			"price": price,
		},
	})
}

// UpdateProduct 更新商品（管理员）
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	// 实际项目中应该验证并更新商品
	name := c.PostForm("name")
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "商品更新成功",
		"data": gin.H{
			"id":    id,
			"name":  name,
			"price": price,
		},
	})
}

// DeleteProduct 删除商品（管理员）
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	// 实际项目中应该验证并删除商品
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "商品 " + id + " 删除成功",
	})
}

// CreateCategory 创建分类（管理员）
func CreateCategory(c *gin.Context) {
	name := c.PostForm("name")
	// 实际项目中应该创建分类
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "分类创建成功",
		"data": gin.H{
			"id":   3,
			"name": name,
		},
	})
}

// UpdateCategory 更新分类（管理员）
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")
	// 实际项目中应该更新分类
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "分类更新成功",
		"data": gin.H{
			"id":   id,
			"name": name,
		},
	})
}

// DeleteCategory 删除分类（管理员）
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	// 实际项目中应该删除分类
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "分类 " + id + " 删除成功",
	})
}
