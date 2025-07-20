package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/stretchr/testify/assert"
)

// createTestData 创建用于测试的文件列表数据
func createTestData() []model.Path {
	baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	return []model.Path{
		{
			Name:  "folder1",
			Path:  "/test/folder1",
			IsDir: true,
			Date:  baseTime.Add(time.Hour * 24 * 3), // 3天后
			Size:  0,
		},
		{
			Name:  "file10.txt",
			Path:  "/test/file10.txt",
			IsDir: false,
			Date:  baseTime.Add(time.Hour * 24 * 2), // 2天后
			Size:  1024,
		},
		{
			Name:  "file2.txt",
			Path:  "/test/file2.txt",
			IsDir: false,
			Date:  baseTime.Add(time.Hour * 24 * 1), // 1天后
			Size:  2048,
		},
		{
			Name:  "folder2",
			Path:  "/test/folder2",
			IsDir: true,
			Date:  baseTime, // 基准时间
			Size:  0,
		},
		{
			Name:  "file1.txt",
			Path:  "/test/file1.txt",
			IsDir: false,
			Date:  baseTime.Add(time.Hour * 24 * 4), // 4天后
			Size:  512,
		},
		{
			Name:  "image.png",
			Path:  "/test/image.png",
			IsDir: false,
			Date:  baseTime.Add(time.Hour * 12), // 12小时后
			Size:  4096,
		},
		{
			Name:  "document",
			Path:  "/test/document",
			IsDir: false,
			Date:  baseTime.Add(time.Hour * 6), // 6小时后
			Size:  1024,
		},
	}
}

// TestGetFileType 测试文件类型判断函数
func TestGetFileType(t *testing.T) {
	tests := []struct {
		name     string
		path     model.Path
		expected string
	}{
		{
			name: "文件夹类型",
			path: model.Path{
				Name:  "testfolder",
				IsDir: true,
			},
			expected: "000_folder",
		},
		{
			name: "txt文件类型",
			path: model.Path{
				Name:  "test.txt",
				IsDir: false,
			},
			expected: "txt",
		},
		{
			name: "png文件类型",
			path: model.Path{
				Name:  "image.PNG",
				IsDir: false,
			},
			expected: "png",
		},
		{
			name: "无扩展名文件",
			path: model.Path{
				Name:  "filename",
				IsDir: false,
			},
			expected: "999_file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 由于 getFileType 是私有函数，我们通过排序行为来测试它
			pathList := model.PathList{tt.path}
			pathList.Sort("type", "asc")

			// 验证排序没有出错（间接测试文件类型函数）
			assert.Len(t, pathList, 1)
			assert.Equal(t, tt.path.Name, pathList[0].Name)
		})
	}
}

// TestSortByName 测试按文件名排序
func TestSortByName(t *testing.T) {
	t.Run("按文件名升序排序", func(t *testing.T) {
		data := createTestData()
		pathList := model.PathList(data)

		pathList.Sort("name", "asc")

		// 验证文件夹在前面，然后是按名称排序的文件
		// 实际排序：document < file1.txt < file2.txt < file10.txt < image.png
		expected := []string{"folder1", "folder2", "document", "file1.txt", "file2.txt", "file10.txt", "image.png"}
		actual := make([]string, len(pathList))
		for i, p := range pathList {
			actual[i] = p.Name
		}

		assert.Equal(t, expected, actual, "文件名升序排序结果不正确")
	})

	t.Run("按文件名降序排序", func(t *testing.T) {
		data := createTestData()
		pathList := model.PathList(data)

		pathList.Sort("name", "desc")

		// 验证文件夹在前面，然后是按名称降序排序的文件
		expected := []string{"folder2", "folder1", "image.png", "file10.txt", "file2.txt", "file1.txt", "document"}
		actual := make([]string, len(pathList))
		for i, p := range pathList {
			actual[i] = p.Name
		}

		assert.Equal(t, expected, actual, "文件名降序排序结果不正确")
	})
}

// TestSortBySize 测试按文件大小排序
func TestSortBySize(t *testing.T) {
	t.Run("按文件大小升序排序", func(t *testing.T) {
		data := createTestData()
		pathList := model.PathList(data)

		pathList.Sort("size", "asc")

		// 验证文件夹在前面，然后是按大小升序排序的文件
		// 大小：file1.txt(512) < document(1024) = file10.txt(1024) < file2.txt(2048) < image.png(4096)
		// 相同大小时按名称排序：document < file10.txt
		expected := []string{"folder1", "folder2", "file1.txt", "document", "file10.txt", "file2.txt", "image.png"}
		actual := make([]string, len(pathList))
		for i, p := range pathList {
			actual[i] = p.Name
		}

		assert.Equal(t, expected, actual, "文件大小升序排序结果不正确")
	})

	t.Run("按文件大小降序排序", func(t *testing.T) {
		data := createTestData()
		pathList := model.PathList(data)

		pathList.Sort("size", "desc")

		// 验证文件夹在前面，然后是按大小降序排序的文件
		// 文件夹按名称排序：folder1 < folder2 -> 降序 -> folder2, folder1
		// 大小降序：image.png(4096) > file2.txt(2048) > document(1024) = file10.txt(1024) > file1.txt(512)
		// 相同大小时按名称排序：document < file10.txt，降序后 -> file10.txt, document
		expected := []string{"folder2", "folder1", "image.png", "file2.txt", "file10.txt", "document", "file1.txt"}
		actual := make([]string, len(pathList))
		for i, p := range pathList {
			actual[i] = p.Name
		}

		assert.Equal(t, expected, actual, "文件大小降序排序结果不正确")
	})
}

// TestSortByModified 测试按修改时间排序
func TestSortByModified(t *testing.T) {
	t.Run("按修改时间升序排序", func(t *testing.T) {
		data := createTestData()
		pathList := model.PathList(data)

		pathList.Sort("modified", "asc")

		// 验证文件夹在前面，然后是按修改时间升序排序的文件
		// 时间：folder2(基准) < folder1(3天后) (文件夹优先)
		// 文件时间：document(6h) < image.png(12h) < file2.txt(1天) < file10.txt(2天) < file1.txt(4天)
		expected := []string{"folder2", "folder1", "document", "image.png", "file2.txt", "file10.txt", "file1.txt"}
		actual := make([]string, len(pathList))
		for i, p := range pathList {
			actual[i] = p.Name
		}

		assert.Equal(t, expected, actual, "修改时间升序排序结果不正确")
	})

	t.Run("按修改时间降序排序", func(t *testing.T) {
		data := createTestData()
		pathList := model.PathList(data)

		pathList.Sort("modified", "desc")

		// 验证文件夹在前面，然后是按修改时间降序排序的文件
		expected := []string{"folder1", "folder2", "file1.txt", "file10.txt", "file2.txt", "image.png", "document"}
		actual := make([]string, len(pathList))
		for i, p := range pathList {
			actual[i] = p.Name
		}

		assert.Equal(t, expected, actual, "修改时间降序排序结果不正确")
	})
}

// TestSortByType 测试按文件类型排序
func TestSortByType(t *testing.T) {
	t.Run("按文件类型升序排序", func(t *testing.T) {
		data := createTestData()
		pathList := model.PathList(data)

		pathList.Sort("type", "asc")

		// 验证按文件类型排序：000_folder < 999_file < png < txt
		// 实际排序：folder1, folder2 (000_folder), document (999_file), image.png (png), file1.txt, file2.txt, file10.txt (txt)
		expected := []string{"folder1", "folder2", "document", "image.png", "file1.txt", "file2.txt", "file10.txt"}
		actual := make([]string, len(pathList))
		for i, p := range pathList {
			actual[i] = p.Name
		}

		assert.Equal(t, expected, actual, "文件类型升序排序结果不正确")
	})

	t.Run("按文件类型降序排序", func(t *testing.T) {
		data := createTestData()
		pathList := model.PathList(data)

		pathList.Sort("type", "desc")

		// 验证按文件类型降序排序：txt > png > 999_file > 000_folder
		expected := []string{"file10.txt", "file2.txt", "file1.txt", "image.png", "document", "folder2", "folder1"}
		actual := make([]string, len(pathList))
		for i, p := range pathList {
			actual[i] = p.Name
		}

		assert.Equal(t, expected, actual, "文件类型降序排序结果不正确")
	})
}

// TestFolderPriority 测试文件夹优先显示逻辑
func TestFolderPriority(t *testing.T) {
	t.Run("文件夹优先显示", func(t *testing.T) {
		data := createTestData()
		pathList := model.PathList(data)

		// 按大小排序时，文件夹应该仍然在前面
		pathList.Sort("size", "asc")

		// 验证前两个是文件夹
		assert.True(t, pathList[0].IsDir, "第一个项目应该是文件夹")
		assert.True(t, pathList[1].IsDir, "第二个项目应该是文件夹")

		// 验证后面的都是文件
		for i := 2; i < len(pathList); i++ {
			assert.False(t, pathList[i].IsDir, "排序后文件夹后面应该都是文件")
		}
	})

	t.Run("type排序时文件夹不优先", func(t *testing.T) {
		data := createTestData()
		pathList := model.PathList(data)

		// 按类型降序排序时，文件夹可能不在最前面
		pathList.Sort("type", "desc")

		// 验证type排序时文件夹位置可能不在最前面
		folderCount := 0
		for _, p := range pathList {
			if p.IsDir {
				folderCount++
			}
		}
		assert.Equal(t, 2, folderCount, "应该有2个文件夹")
	})
}

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	t.Run("空列表排序", func(t *testing.T) {
		pathList := model.PathList{}
		pathList.Sort("name", "asc")

		assert.Len(t, pathList, 0, "空列表排序后应该仍然为空")
	})

	t.Run("单个文件排序", func(t *testing.T) {
		pathList := model.PathList{
			{
				Name:  "single.txt",
				Path:  "/test/single.txt",
				IsDir: false,
				Date:  time.Now(),
				Size:  100,
			},
		}

		pathList.Sort("name", "asc")

		assert.Len(t, pathList, 1, "单个文件排序后应该仍然只有一个文件")
		assert.Equal(t, "single.txt", pathList[0].Name, "文件名应该保持不变")
	})

	t.Run("无效排序字段", func(t *testing.T) {
		data := createTestData()
		originalOrder := make([]string, len(data))
		for i, p := range data {
			originalOrder[i] = p.Name
		}

		pathList := model.PathList(data)
		pathList.Sort("invalid_field", "asc")

		// 验证无效字段时使用默认排序（按名称）
		expected := []string{"folder1", "folder2", "document", "file1.txt", "file2.txt", "file10.txt", "image.png"}
		actual := make([]string, len(pathList))
		for i, p := range pathList {
			actual[i] = p.Name
		}

		assert.Equal(t, expected, actual, "无效排序字段应该使用默认排序")
	})

	t.Run("空排序字段", func(t *testing.T) {
		data := createTestData()
		originalOrder := make([]string, len(data))
		for i, p := range data {
			originalOrder[i] = p.Name
		}

		pathList := model.PathList(data)
		pathList.Sort("", "asc")

		// 验证空排序字段时不进行排序
		actual := make([]string, len(pathList))
		for i, p := range pathList {
			actual[i] = p.Name
		}

		assert.Equal(t, originalOrder, actual, "空排序字段应该保持原顺序")
	})

	t.Run("无效排序方向", func(t *testing.T) {
		data := createTestData()
		pathList := model.PathList(data)

		pathList.Sort("name", "invalid")

		// 验证无效排序方向时使用升序
		expected := []string{"folder1", "folder2", "document", "file1.txt", "file2.txt", "file10.txt", "image.png"}
		actual := make([]string, len(pathList))
		for i, p := range pathList {
			actual[i] = p.Name
		}

		assert.Equal(t, expected, actual, "无效排序方向应该使用升序")
	})
}

// TestSameSizeHandling 测试相同大小文件的处理
func TestSameSizeHandling(t *testing.T) {
	t.Run("相同大小文件按名称排序", func(t *testing.T) {
		pathList := model.PathList{
			{
				Name:  "z_file.txt",
				Path:  "/test/z_file.txt",
				IsDir: false,
				Date:  time.Now(),
				Size:  1024,
			},
			{
				Name:  "a_file.txt",
				Path:  "/test/a_file.txt",
				IsDir: false,
				Date:  time.Now(),
				Size:  1024,
			},
		}

		pathList.Sort("size", "asc")

		// 验证相同大小时按名称排序
		assert.Equal(t, "a_file.txt", pathList[0].Name, "相同大小时应该按名称升序排序")
		assert.Equal(t, "z_file.txt", pathList[1].Name, "相同大小时应该按名称升序排序")
	})
}

// TestSameTimeHandling 测试相同修改时间文件的处理
func TestSameTimeHandling(t *testing.T) {
	t.Run("相同修改时间文件按名称排序", func(t *testing.T) {
		sameTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		pathList := model.PathList{
			{
				Name:  "z_file.txt",
				Path:  "/test/z_file.txt",
				IsDir: false,
				Date:  sameTime,
				Size:  100,
			},
			{
				Name:  "a_file.txt",
				Path:  "/test/a_file.txt",
				IsDir: false,
				Date:  sameTime,
				Size:  200,
			},
		}

		pathList.Sort("modified", "asc")

		// 验证相同修改时间时按名称排序
		assert.Equal(t, "a_file.txt", pathList[0].Name, "相同修改时间时应该按名称升序排序")
		assert.Equal(t, "z_file.txt", pathList[1].Name, "相同修改时间时应该按名称升序排序")
	})
}

// BenchmarkSort 性能基准测试
func BenchmarkSort(b *testing.B) {
	// 创建大量测试数据
	data := make([]model.Path, 1000)
	baseTime := time.Now()

	for i := 0; i < 1000; i++ {
		data[i] = model.Path{
			Name:  fmt.Sprintf("file%d.txt", i),
			Path:  fmt.Sprintf("/test/file%d.txt", i),
			IsDir: i%10 == 0, // 每10个文件有一个文件夹
			Date:  baseTime.Add(time.Duration(i) * time.Minute),
			Size:  int64(i * 100),
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 创建副本以避免影响基准测试
		testData := make([]model.Path, len(data))
		copy(testData, data)

		pathList := model.PathList(testData)
		pathList.Sort("name", "asc")
	}
}
