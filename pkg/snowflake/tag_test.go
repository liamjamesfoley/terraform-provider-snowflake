package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagCreate(t *testing.T) {
	r := require.New(t)
	o := Tag("test")
	r.Equal(o.QualifiedName(), `"test"`)

	o.WithDB("db")
	r.Equal(o.QualifiedName(), `"db"."test"`)

	o.WithSchema("schema")
	r.Equal(o.QualifiedName(), `"db"."schema"."test"`)

	r.Equal(o.Create(), `CREATE TAG "db"."schema"."test"`)

	allowedValues := []string{"marketing", "finance"}
	o.WithAllowedValues(allowedValues)

	o.WithComment("Yee'haw")
	r.Equal(`CREATE TAG "db"."schema"."test" ALLOWED_VALUES 'marketing', 'finance' COMMENT = 'Yee\'haw'`, o.Create())
}

func TestTagRename(t *testing.T) {
	r := require.New(t)
	o := Tag("test")
	r.Equal(o.Rename("bob"), `ALTER TAG "test" RENAME TO "bob"`)
}

func TestTagChangeComment(t *testing.T) {
	r := require.New(t)
	o := Tag("test")
	r.Equal(`ALTER TAG "test" SET COMMENT = 'worst\' tag ever'`, o.ChangeComment("worst' tag ever"))
}

func TestTagRemoveComment(t *testing.T) {
	r := require.New(t)
	o := Tag("test")
	r.Equal(o.RemoveComment(), `ALTER TAG "test" UNSET COMMENT`)
}

func TestTagAddAllowedValues(t *testing.T) {
	r := require.New(t)
	o := Tag("test")
	avs := []string{"foo", "bar"}
	r.Equal(o.AddAllowedValues(avs), `ALTER TAG "test" ADD ALLOWED_VALUES 'foo', 'bar'`)
}

func TestTagDropAllowedValues(t *testing.T) {
	r := require.New(t)
	o := Tag("test")
	avs := []string{"foo"}
	r.Equal(o.DropAllowedValues(avs), `ALTER TAG "test" DROP ALLOWED_VALUES 'foo'`)
}

func TestTagRemoveAllowedValues(t *testing.T) {
	r := require.New(t)
	o := Tag("test")
	r.Equal(o.RemoveAllowedValues(), `ALTER TAG "test" UNSET ALLOWED_VALUES`)
}

func TestTagDrop(t *testing.T) {
	r := require.New(t)
	o := Tag("test")
	r.Equal(o.Drop(), `DROP TAG "test"`)
}

func TestTagUndrop(t *testing.T) {
	r := require.New(t)
	o := Tag("test")
	r.Equal(o.Undrop(), `UNDROP TAG "test"`)
}

func TestTagShow(t *testing.T) {
	r := require.New(t)
	o := Tag("test")
	r.Equal(o.Show(), `SHOW TAGS LIKE 'test'`)

	o.WithDB("db")
	r.Equal(o.Show(), `SHOW TAGS LIKE 'test' IN DATABASE "db"`)

	o.WithSchema("schema")
	r.Equal(o.Show(), `SHOW TAGS LIKE 'test' IN SCHEMA "db"."schema"`)
}