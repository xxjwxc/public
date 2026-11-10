package fractional

import "testing"

func TestModel(t *testing.T) {
	if f := Model(3); f.Nume != 3 || f.Deno != 1 {
		t.Errorf("Model(3) 期望 3/1,实际 %d/%d", f.Nume, f.Deno)
	}
	if f := Model(3, 4); f.Nume != 3 || f.Deno != 4 {
		t.Errorf("Model(3,4) 期望 3/4,实际 %d/%d", f.Nume, f.Deno)
	}

	defer func() {
		if recover() == nil {
			t.Error("Model(1,0) 应 panic")
		}
	}()
	Model(1, 0)
}

func TestAdd(t *testing.T) {
	tmp := Model(7, 12)
	tmp1 := Model(1, 12)

	got := tmp.Add(tmp1)
	// 7/12 + 1/12 = 8/12 = 2/3
	if got.Nume != 2 || got.Deno != 3 {
		t.Fatalf("7/12+1/12 期望 2/3,实际 %d/%d", got.Nume, got.Deno)
	}

	// 入参不应被修改
	if tmp1.Nume != 1 || tmp1.Deno != 12 {
		t.Errorf("Add 修改了入参,现在为 %d/%d,期望 1/12", tmp1.Nume, tmp1.Deno)
	}
}

func TestSub(t *testing.T) {
	tmp := Model(1, 4)
	tmp1 := Model(1, 3)

	got := tmp.Sub(tmp1)
	// 1/4 - 1/3 = -1/12
	if got.Nume != -1 || got.Deno != 12 {
		t.Fatalf("1/4-1/3 期望 -1/12,实际 %d/%d", got.Nume, got.Deno)
	}
	if tmp1.Nume != 1 || tmp1.Deno != 3 {
		t.Errorf("Sub 修改了入参,现在为 %d/%d,期望 1/3", tmp1.Nume, tmp1.Deno)
	}
}

func TestMulDiv(t *testing.T) {
	if got := Model(3, 4).Mul(Model(2, 3)); got.Nume != 1 || got.Deno != 2 {
		t.Errorf("3/4*2/3 期望 1/2,实际 %d/%d", got.Nume, got.Deno)
	}
	// (3/4) / (2/3) = 9/8
	if got := Model(3, 4).Div(Model(2, 3)); got.Nume != 9 || got.Deno != 8 {
		t.Errorf("(3/4)/(2/3) 期望 9/8,实际 %d/%d", got.Nume, got.Deno)
	}
}

func TestVerdict(t *testing.T) {
	if got := Model(1, 3).Verdict(); got < 0.333 || got > 0.334 {
		t.Errorf("1/3 的值期望约 0.3333,实际 %v", got)
	}
}
