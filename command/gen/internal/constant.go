package internal

const (
	// OverwriteTag header
	OverwriteTag = "// this file will be overwritten when execute gen command next."
	TemplateDir  = "template"
	// EmptyData DataType
	EmptyData     = "empty"     // 空代码文件
	ModelData     = "model"     // 基于单个模型（表结构）生成代码
	GeneratorData = "generator" // 基于“代码生成器”配置数据生成代码
)

// 数据库-数据类型
const (
	Text       = "text"          // 文本
	Varchar    = "varchar"       // 字符
	Varchar100 = "varchar(100)"  // 100字符串
	Char       = "char"          // 字节
	Int        = "int"           // 整型
	Int2       = "int2"          // 小整型
	Int4       = "int4"          // 中整型
	Int8       = "int8"          // 大整型
	Tinyint    = "tinyint"       // 微整数
	Smallint   = "smallint"      // 小整数
	Mediumint  = "mediumint"     // 中整数
	Bigint     = "bigint"        // 大整数
	Float      = "float"         // 浮点数
	Float4     = "float4"        // 浮点数
	Float8     = "float8"        // 浮点数
	Decimal    = "decimal"       // 十进制数
	Sequence   = "sequence"      // 序列
	Bool       = "bool"          // 布尔
	Uuid       = "uuid"          // UUID
	Numeric    = "numeric"       // 数字
	Numeric2   = "numeric(10,2)" // 2精度数字
	Time       = "time"          // 时间
	Date       = "date"          // 日期
	Timestamp  = "timestamp"     // 时间戳
	Timestampz = "timestamptz"   // 时间戳带时区
	Datetime   = "datetime"      // 日期时间
)

// go基础数据类型
const (
	GoString  = "string"
	GoInt     = "int"
	GoInt64   = "int64"
	GoFloat64 = "float64"
	GoBool    = "bool"
	GoTime    = "time.Time"
)

// Java基础数据类型
const (
	JavaString     = "String"
	JavaInteger    = "Integer"
	JavaInt        = "int"
	JavaLong       = "Long"
	JavaDate       = "Date"
	JavaBigDecimal = "BigDecimal"
	JavaFloat      = "Float"
	JavaBoolean    = "Boolean"
)

// IsBaseField BASE基础字段映射
func IsBaseField(t string) bool {
	switch t {
	case "id":
		return true
	case "create_time", "create_user_id", "create_by":
		return true
	case "update_time", "update_user_id", "update_by":
		return true
	default:
		return false
	}
}

// DB2GoType DB-Go类型映射
func DB2GoType(t string) string {
	switch t {
	case Char, Varchar, Varchar100, Text, Uuid:
		return GoString
	case Int, Int2, Int4, Tinyint, Smallint, Mediumint:
		return GoInt
	case Int8, Bigint:
		return GoInt64
	case Float, Float4, Float8, Numeric:
		return GoFloat64
	case Timestamp, Timestampz, Datetime, Time, Date:
		return GoTime
	case Bool:
		return GoBool
	default:
		return GoString
	}
}

// DB2GormType DB-Gorm类型映射
func DB2GormType(t string) string {
	switch t {
	case Char:
		return Char
	case Varchar:
		return Varchar100
	case Text:
		return Text
	case Tinyint:
		return Tinyint
	case Smallint, Int2:
		return Smallint
	case Mediumint, Int4, Int:
		return Int
	case Bigint, Int8:
		return Bigint
	case Float4, Numeric:
		return Numeric2
	case Timestamp, Timestampz, Datetime:
		return Timestamp
	case Bool:
		return Bool
	case Date:
		return Date
	default:
		return t
	}
}

// DB2JavaType DB-JAVA类型映射
func DB2JavaType(t string) string {
	switch t {
	case Char, Varchar, Text:
		return JavaString
	case Int2, Smallint, Mediumint:
		return JavaInteger
	case Tinyint:
		return JavaInt
	case Int, Int4, Int8, Bigint:
		return JavaLong
	case Float4, Numeric:
		return JavaFloat
	case Decimal:
		return JavaBigDecimal
	case Timestamp, Datetime, Date:
		return JavaDate
	case Bool:
		return JavaBoolean
	default:
		return JavaString
	}
}
