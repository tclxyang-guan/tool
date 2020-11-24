/**
* @Auther:gy
* @Date:2020/11/23 16:40
 */

package models

//自定义gorm Model
type Model struct {
	ID        uint    `gorm:"primary_key;comment:'数据id'" json:"id" req:"-"`
	CreatedAt string  `gorm:"type:varchar(30);comment:'创建时间'" json:"created_at,omitempty" req:"-"`
	UpdatedAt string  `gorm:"type:varchar(30);comment:'修改时间'" json:"updated_at,omitempty" req:"-"`
	DeletedAt *string `gorm:"type:varchar(30);default:null;comment:'删除时间'" json:"deleted_at" req:"-" resp:"-"`
}
type Build struct {
	Model
	BuildName      string `gorm:"default:'';comment:'楼栋名称'" json:"build_name" validate:"required"`
	OrganizationID uint   `gorm:"default:0;index:org_com_vil;not null;comment:'组织id'" json:"organization_id" validate:"required"`
	CommunityID    uint   `gorm:"default:0;index:org_com_vil,com_status;not null;comment:'社区id'" json:"community_id" validate:"required"`
	NeighborID     uint   `gorm:"default:0;index:org_com_vil,vil;not null;comment:'小区id'" json:"neighbor_id"`
	Address        string `gorm:"default:'';comment:'地址'" json:"address" validate:"required"`
	Lon            string `gorm:"default:'';comment:'经度'" json:"lon" validate:"required"`
	Lat            string `gorm:"default:'';comment:'纬度'" json:"lat" validate:"required"`
	DetailAddress  string `gorm:"not null;comment:'详细地址'"json:"detail_address" validate:"required"`
	UpCount        uint   `gorm:"default:0;comment:'地上层数'" json:"up_count" validate:"required"`
	DownCount      uint   `gorm:"default:0;comment:'地下层数'" json:"down_count" validate:"required"`
	Unit           string `gorm:"default:0;comment:'所属单元'" json:"unit"`
	Count          uint   `gorm:"default:0;comment:'总户数'" json:"count" validate:"required"`
	Lift           uint   `gorm:"default:0;comment:'电梯数'" json:"lift" validate:"required"`
	Status         uint   `gorm:"default:2;index:com_status;comment:'发布状态 1已发布 2待发布 '"json:"status" req:"-"`
}
