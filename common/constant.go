package common

const (
	String   = "string"
	Int      = "int"      // 数字
	Float    = "float"    // 浮点数
	Sequence = "sequence" // 数字编号
	Time     = "time"     // 时间
	Date     = "date"     // 日期
	Uuid     = "uuid"     // UUID
	Phone    = "phone"    // 手机号
	Name     = "name"     // 姓名
	IdCard   = "id_card"  // 身份证
	PlateNo  = "plate_no" // 车牌号
	Email    = "email"    // 邮箱
	IP       = "ip"       // ip地址
	Province = "province" // 省
	City     = "city"     // 市
	Password = "password" // 密码
	Enum     = "enum"     // 枚举
	Database = "database" // 数据库

	Int2      = "int2"
	Int4      = "int4"
	Int8      = "int8"
	Tinyint   = "tinyint"
	Smallint  = "smallint"
	Mediumint = "mediumint"
	Bigint    = "bigint"
	Float4    = "float4"
	Numeric   = "numeric" // 数字
	Timestamp = "timestamp"
	TimeUnix  = "unix"
	Datetime  = "datetime"

	Backup  = "backup"  // restic备份
	Restore = "restore" // restic恢复
	Forget  = "forget"  // restic删除
	Init    = "init"    // restic初始化
)

// DB2RandType DB-Go类型映射
func DB2RandType(t string) string {
	switch t {
	case Int, Int2, Int4, Int8, Tinyint, Smallint, Mediumint, Bigint:
		return Int
	case Float4, Numeric:
		return Float
	case Timestamp, Datetime:
		return Time
	case Date:
		return Date
	default:
		return String
	}
}
