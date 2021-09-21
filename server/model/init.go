package model

import (
	"gin-server-skeleton/db"
	"os"
	"time"
)

// init
func init() {
	if runMode := os.Getenv("RUN_MODE"); runMode == "testing" {
		// TO-DO
	} else {
		AutoMigrateAll()
	}
}

// Migrate Model
func AutoMigrateAll() {
	_ = db.Conn.Table("users").AutoMigrate(&User{})
	//db.Conn.Exec("INSERT INTO gin_scaffold.users (`id`,`name`,`password`,`email`,`age`,`birthday`,`member_number`,`created_at`,`updated_at`,`deleted_at`,`role`) VALUES (1,'admin','$2a$10$BjYFeoOSaD8Xzs2KumA7Z.duVszjK8lB1ZaZDkdlc5bTzvPSbGify','admin@admin.com',0,'2021-05-24 23:39:01.278','1','2021-05-24 23:39:01.280','2021-05-24 23:39:01.280',NULL,'管理员'); ")
	admin := User{
		ID:       1,
		Name:     "admin",
		Password: "$2a$10$BjYFeoOSaD8Xzs2KumA7Z.duVszjK8lB1ZaZDkdlc5bTzvPSbGify",
		Email:    "admin@admin.com",
		Age:      0,
		Birthday: time.Now(),
		Role:     "管理员",
		BaseModel: BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	db.Conn.Model(&User{}).Create(admin)
}
