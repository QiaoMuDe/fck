# 注册表重构设计方案

## 问题背景

当前注册表实现存在以下问题: 
1. 标志对象被重复注册 (长名称和短名称各注册一次) 
2. 内置标志 (如 help) 与用户自定义标志可能冲突
3. 测试中期望的标志数量与实际不符
4. 命令注册表和标志注册表的处理逻辑不一致

## 设计目标

1. 避免重复存储同一个标志对象
2. 统一命令注册表和标志注册表的处理逻辑
3. 保持与底层标准库的兼容性
4. 提供高效的查找机制
5. 简化代码逻辑, 提高可维护性

## 方案设计

### 数据结构

```go
type registry[T any] struct {
	items     map[int]T         // 主存储, 使用数字索引存储实际对象 (只存储一次) 
	nameIndex map[string]int     // 名称到索引的映射 (支持长名称和短名称, 都指向同一个索引) 
	nextID    int              // 下一个可用的ID (自增) 
}
```

### 核心方法

#### NewRegistry
创建新的注册表实例
```go
func NewRegistry[T any]() *registry[T] {
	return &registry[T]{
		items:     make(map[int]T),
		nameIndex: make(map[string]int),
		nextID:    1, // 从1开始, 0表示无效ID
	}
}
```

#### Register
注册项到注册表, 支持长名称和短名称
```go
func (r *registry[T]) Register(item T, longName, shortName string) error {
	// 第一步: 验证参数
	if longName == "" && shortName == "" {
		return types.NewError("INVALID_NAME", "long name and short name cannot both be empty", nil)
	}
	
	// 第二步: 检查名称冲突
	if longName != "" {
		if _, exists := r.nameIndex[longName]; exists {
			return types.ErrAlreadyExists
		}
	}
	
	if shortName != "" {
		if _, exists := r.nameIndex[shortName]; exists {
			return types.ErrAlreadyExists
		}
	}
	
	// 第三步: 获取索引并存储到主存储
	id := r.nextID
	r.nextID++
	r.items[id] = item  // 直接获取索引, 把实际对象存储进去
	
	// 第四步: 建立名称索引映射
	// 通过判断长短名称如果不为空就注册一个, value就是实际的索引
	if longName != "" {
		r.nameIndex[longName] = id  // 长名称映射到索引
	}
	
	if shortName != "" {
		r.nameIndex[shortName] = id  // 短名称映射到索引
	}
	
	return nil
}
```

#### Get
通过名称获取项 (支持长名称或短名称) 
```go
func (r *registry[T]) Get(name string) (T, bool) {
	id, exists := r.nameIndex[name]
	if !exists {
		var zero T
		return zero, false
	}
	
	item, exists := r.items[id]
	return item, exists
}
```

#### GetByShortName
通过短名称获取项
```go
func (r *registry[T]) GetByShortName(shortName string) (T, bool) {
	// 直接调用Get方法, 因为内部已经统一处理
	return r.Get(shortName)
}
```

#### Unregister
通过名称注销项
```go
func (r *registry[T]) Unregister(name string) error {
	id, exists := r.nameIndex[name]
	if !exists {
		return types.ErrNotFound
	}
	
	// 找出所有指向该ID的名称并删除
	for name, nameID := range r.nameIndex {
		if nameID == id {
			delete(r.nameIndex, name)
		}
	}
	
	// 删除项
	delete(r.items, id)
	
	return nil
}
```

#### 其他辅助方法
```go
// List 获取所有项
func (r *registry[T]) List() []T {
	result := make([]T, 0, len(r.items))
	for _, item := range r.items {
		result = append(result, item)
	}
	return result
}

// Count 获取项的数量
func (r *registry[T]) Count() int {
	return len(r.items)
}

// Has 检查名称是否存在
func (r *registry[T]) Has(name string) bool {
	_, exists := r.nameIndex[name]
	return exists
}

// Clear 清空注册表
func (r *registry[T]) Clear() {
	r.items = make(map[int]T)
	r.nameIndex = make(map[string]int)
	r.nextID = 1
}

// Range 遍历所有项
func (r *registry[T]) Range(fn func(T) bool) {
	for _, item := range r.items {
		if !fn(item) {
			break
		}
	}
}
```

### 适配现有接口

#### FlagRegistryImpl适配
```go
// 为FlagRegistryImpl提供适配
func (r *FlagRegistryImpl) Register(flag types.Flag) error {
	return r.registry.Register(flag, flag.Name(), flag.ShortName())
}
```

#### CmdRegistryImpl适配
```go
// 为CmdRegistryImpl提供适配
func (r *CmdRegistryImpl) Register(cmd types.Command) error {
	return r.registry.Register(cmd, cmd.Name(), cmd.ShortName())
}
```

## 方案优势

1. **避免重复存储**: 标志对象只存储一次, 长名称和短名称都指向同一个对象
2. **统一处理**: 命令和标志使用相同的底层逻辑
3. **高效查找**: 通过名称直接映射到ID, 再通过ID获取对象
4. **内存效率**: 相比当前方案, 减少了重复对象的存储
5. **简洁设计**: 去掉了不必要的锁机制, 简化了代码逻辑

## 与解析器的协作

### 解析器层面
解析器仍然需要分别注册长名称和短名称到标准库的flagSet中, 确保命令行参数能正确解析: 
```go
func (p *DefaultParser) registerFlag(f types.Flag) {
	longName := f.LongName()
	shortName := f.ShortName()
	description := f.Desc()

	// 注册长名称
	if longName != "" {
		p.flagSet.Var(newFlagValueWrapper(f), longName, description)
	}

	// 注册短名称
	if shortName != "" {
		p.flagSet.Var(newFlagValueWrapper(f), shortName, description)
	}
}
```

### 注册表层面
注册表只存储一次标志对象, 但支持通过长名称或短名称查找, 这样既避免了重复存储, 又保持了查找的灵活性。

## 实施计划

1. 重构 `registry/impl.go` 中的数据结构和方法
2. 更新 `registry/flag_registry.go` 中的Register方法
3. 更新 `registry/command_registry.go` 中的Register方法
4. 修改相关测试用例, 适应新的实现
5. 验证所有测试通过

## 注意事项

1. **名称冲突**: 长名称和短名称可能相同, 需要在注册时检查
2. **ID回收**: 删除项后ID不会被重用, 这是可接受的
3. **内存增长**: nameIndex可能比items大, 这是预期的, 因为一个项可能有多个名称
4. **向后兼容**: 确保现有API保持不变, 只改变内部实现