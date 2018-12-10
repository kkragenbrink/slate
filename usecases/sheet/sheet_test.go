// Copyright (c) 2018 Kevin Kragenbrink, II
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package sheet

import (
	"context"
	"github.com/bmizerany/assert"
	"github.com/golang/mock/gomock"
	"github.com/kkragenbrink/slate/domains"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SheetSuite struct {
	suite.Suite
}

func TestSheet(t *testing.T) {
	suite.Run(t, new(SheetSuite))
}

func (suite *SheetSuite) TestNew() {
	ctrl, ctx := gomock.WithContext(context.Background(), suite.T())
	db := domains.NewMockCharacterRepository(ctrl)
	db.EXPECT().Store(ctx, gomock.Any()).Return(nil)
	ch, err := New(ctx, db, "Test", "wtf2e", int64(1234), int64(1234))
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), "wtf2e", ch.Sheet.System())
}

func (suite *SheetSuite) TestGenerateSheetBySystem() {
	sh := GenerateSheetBySystem("wtf2e", nil)
	assert.Equal(suite.T(), "wtf2e", sh.System())
}
