package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompletePasswordPolicy(t *testing.T) {
	r := require.New(t)
	client, err := NewDefaultClient()
	r.Nil(err)
	objectIdentifier := NewSchemaObjectIdentifier("TEST_DB", "PUBLIC", "test_policy")
	name := objectIdentifier.FullyQualifiedName()
	err = client.PasswordPolicies.Create(context.Background(), name, &PasswordPolicyCreateOptions{
		OrReplace:                 Bool(true),
		PasswordMinLength:         Int(10),
		PasswordMaxLength:         Int(20),
		PasswordMinUpperCaseChars: Int(1),
		PasswordMinLowerCaseChars: Int(1),
		PasswordMinNumericChars:   Int(1),
		PasswordMinSpecialChars:   Int(1),
		PasswordMaxAgeDays:        Int(30),
		PasswordMaxRetries:        Int(5),
		PasswordLockoutTimeMins:   Int(30),
		Comment:                   String("test"),
	})
	r.Nil(err)

	_, err = client.PasswordPolicies.Show(context.Background(), &PasswordPolicyShowOptions{
		Pattern: String(objectIdentifier.Name),
		In: &PasswordPolicyShowIn{
			Database: String(objectIdentifier.DatabaseName),
		},
	})
	r.Nil(err)

	err = client.PasswordPolicies.Alter(context.Background(), name, &PasswordPolicyAlterOptions{
		Set: &PasswordPolicyAlterSet{
			PasswordMinLength: Int(8),
			Comment:           String("test22"),
		},
	})
	r.Nil(err)

	_, err = client.PasswordPolicies.Describe(context.Background(), name)
	r.Nil(err)

	err = client.PasswordPolicies.Drop(context.Background(), name, &PasswordPolicyDropOptions{
		IfExists: Bool(true),
	})
	r.Nil(err)
}