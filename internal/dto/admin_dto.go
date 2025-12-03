package dto

// AdminDashboardResp 后台仪表盘响应。
type AdminDashboardResp struct {
	Metrics  AdminDashboardMetrics `json:"metrics"`
	Recent   AdminRecentStats      `json:"recent"`
	TopPosts []AdminTopPost        `json:"top_posts"`
}

// AdminDashboardMetrics 汇总总量指标。
type AdminDashboardMetrics struct {
	Users    int64 `json:"users"`
	Posts    int64 `json:"posts"`
	Comments int64 `json:"comments"`
}

// AdminRecentStats 近7天新增指标。
type AdminRecentStats struct {
	NewUsers    int64 `json:"new_users"`
	NewPosts    int64 `json:"new_posts"`
	NewComments int64 `json:"new_comments"`
}

// AdminTopPost 评论最多的文章。
type AdminTopPost struct {
	PostID       uint   `json:"post_id"`
	Title        string `json:"title"`
	CommentCount int64  `json:"comment_count"`
}

// AdminUserQuery 用户列表查询条件。
type AdminUserQuery struct {
	Page     int
	PageSize int
	Keyword  string
	Role     string
}

// AdminPostQuery 文章列表查询条件。
type AdminPostQuery struct {
	Page     int
	PageSize int
	Keyword  string
	Status   *string
}

// AdminCommentQuery 评论列表查询条件。
type AdminCommentQuery struct {
	Page     int
	PageSize int
	Keyword  string
	UserID   *uint
	PostID   *uint
}
