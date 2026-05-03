# 🤝 Contributing to Omnilog

Дякуємо за інтерес до проєкту! Ми вітаємо будь-які внески - від виправлення друкарських помилок до додавання нових функцій.

## 📋 Зміст

- [Code of Conduct](#code-of-conduct)
- [Як допомогти](#як-допомогти)
- [Процес розробки](#процес-розробки)
- [Coding Guidelines](#coding-guidelines)
- [Commit Convention](#commit-convention)
- [Pull Request Process](#pull-request-process)

---

## Code of Conduct

Цей проєкт дотримується [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). Беручи участь, ви погоджуєтесь дотримуватись цих правил.

---

## Як допомогти

### 🐛 Повідомити про баг

1. Перевірте [Issues](https://github.com/markostrutinsky/diploma/issues), чи не створена вже така проблема
2. Якщо ні, створіть новий Issue з описом:
   - Кроки для відтворення
   - Очікувана поведінка
   - Фактична поведінка
   - Скріншоти (якщо можливо)
   - Версія системи, браузер, ОС

### 💡 Запропонувати нову функцію

1. Створіть Issue з тегом `enhancement`
2. Опишіть проблему, яку вирішує ваша ідея
3. Запропонуйте можливе рішення
4. Дочекайтесь обговорення з мейнтейнерами

### 🔧 Виправити баг або додати функцію

1. Fork репозиторій
2. Створіть нову гілку (`git checkout -b feature/amazing-feature`)
3. Зробіть зміни
4. Commit (`git commit -m 'feat: add amazing feature'`)
5. Push (`git push origin feature/amazing-feature`)
6. Створіть Pull Request

---

## Процес розробки

### 1. Налаштування середовища

```bash
# Клонуйте ваш fork
git clone https://github.com/your-username/diploma.git
cd diploma

# Додайте upstream remote
git remote add upstream https://github.com/markostrutinsky/diploma.git

# Створіть .env файл
cp .env.example .env
# Відредагуйте .env

# Запустіть Docker
docker compose up -d --build
```

### 2. Структура гілок

- `main` - production-ready код
- `dev` - development гілка (основна для PR)
- `feature/*` - нові функції
- `fix/*` - виправлення багів
- `docs/*` - документація
- `refactor/*` - рефакторинг коду

### 3. Робота з upstream

```bash
# Отримайте останні зміни
git fetch upstream

# Оновіть вашу dev гілку
git checkout dev
git merge upstream/dev

# Створіть feature гілку
git checkout -b feature/my-feature
```

---

## Coding Guidelines

### Backend (Go)

```go
// ✅ Добре
func (s *InventoryService) CreateResource(ctx context.Context, req CreateResourceRequest) (*Resource, error) {
    // Validate input
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }

    // Business logic
    resource := &Resource{
        ID:   uuid.New(),
        Name: req.Name,
        // ...
    }

    // Save to database
    if err := s.repo.Create(ctx, resource); err != nil {
        return nil, fmt.Errorf("failed to create resource: %w", err)
    }

    return resource, nil
}

// ❌ Погано
func create(r req) *res {
    resource := &Resource{ID: uuid.New(), Name: r.Name}
    s.repo.Create(resource)
    return resource
}
```

**Правила:**
- Використовуйте `context.Context` для всіх операцій з БД
- Обробляйте помилки з `fmt.Errorf` та `%w`
- Використовуйте `slog` для логування
- Називайте змінні та функції зрозуміло
- Додавайте коментарі до публічних функцій

### Frontend (React + TypeScript)

```typescript
// ✅ Добре
interface ResourceCardProps {
  resource: Resource;
  onEdit: (id: string) => void;
  onDelete: (id: string) => void;
}

export const ResourceCard: React.FC<ResourceCardProps> = ({ 
  resource, 
  onEdit, 
  onDelete 
}) => {
  const handleEdit = useCallback(() => {
    onEdit(resource.id);
  }, [resource.id, onEdit]);

  return (
    <div className="resource-card">
      <h3>{resource.name}</h3>
      <p>Quantity: {resource.quantity}</p>
      <button onClick={handleEdit}>Edit</button>
    </div>
  );
};

// ❌ Погано
export const ResourceCard = (props) => {
  return (
    <div>
      <h3>{props.resource.name}</h3>
      <button onClick={() => props.onEdit(props.resource.id)}>Edit</button>
    </div>
  );
};
```

**Правила:**
- Використовуйте TypeScript типи для всіх props
- Використовуйте functional components та hooks
- Використовуйте `useCallback` для callback функцій
- CSS Modules або BEM для стилів
- Accessible HTML (aria-labels, semantic tags)

### Database

```sql
-- ✅ Добре - з індексами та constraints
CREATE TABLE resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_resources_tenant_id ON resources(tenant_id);
CREATE INDEX idx_resources_category_id ON resources(category_id);

-- ❌ Погано - без constraints
CREATE TABLE resources (
    id TEXT,
    name TEXT,
    quantity INTEGER
);
```

**Правила:**
- Завжди використовуйте UUID для ID
- Додавайте `tenant_id` для multi-tenant ізоляції
- Використовуйте `TIMESTAMPTZ` для timestamps
- Додавайте індекси для foreign keys та часто запитуваних полів
- Використовуйте CHECK constraints для валідації

---

## Commit Convention

Ми використовуємо [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: Нова функція
- `fix`: Виправлення бага
- `docs`: Зміни в документації
- `style`: Форматування, пропущені крапки з комою, тощо
- `refactor`: Рефакторинг коду
- `perf`: Покращення продуктивності
- `test`: Додавання або виправлення тестів
- `chore`: Оновлення build tasks, package manager configs, тощо

### Examples

```bash
feat(inventory): add Excel import functionality

- Add upload endpoint with multipart/form-data support
- Parse Excel file with SheetJS
- Validate rows and create resources in batch
- Return detailed import report with errors

Closes #123

---

fix(auth): prevent JWT token expiration race condition

The refresh token mechanism had a race condition where multiple
simultaneous requests could cause token invalidation.

Fixed by implementing token refresh lock with Redis.

Fixes #456

---

docs(api): add GPS tracking endpoints to API documentation

- Add /api/gps/* endpoints
- Include request/response examples
- Document geofencing alerts

---

refactor(services): extract common repository logic

- Create base repository interface
- Implement tenant isolation in base repo
- Update all repositories to use base

---

perf(database): add indexes for frequently queried fields

Added indexes:
- resources(tenant_id, category_id)
- requests(tenant_id, status, created_at)
- vehicles(tenant_id, status)

Query time reduced from 250ms to 15ms on large datasets.
```

---

## Pull Request Process

### 1. Перед створенням PR

- [ ] Код компілюється без помилок
- [ ] Немає lint warnings
- [ ] Додані/оновлені тести (якщо потрібно)
- [ ] Документація оновлена (якщо потрібно)
- [ ] Commit messages відповідають Conventional Commits
- [ ] PR створений з `dev` гілки

### 2. Опис PR

Використовуйте шаблон:

```markdown
## Що змінено

Короткий опис змін.

## Тип змін

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Як тестувати

1. Крок 1
2. Крок 2
3. Очікуваний результат

## Скріншоти (якщо UI зміни)

![Screenshot](url)

## Checklist

- [ ] Код відповідає стилю проєкту
- [ ] Зміни протестовані
- [ ] Документація оновлена
- [ ] Немає нових warnings
```

### 3. Review Process

1. Мінімум 1 схвалення від мейнтейнера
2. Всі CI checks пройшли
3. Конфлікти вирішені
4. PR буде змержено в `dev`

### 4. Після merge

```bash
# Оновіть ваш fork
git checkout dev
git pull upstream dev
git push origin dev

# Видаліть feature гілку
git branch -d feature/my-feature
git push origin --delete feature/my-feature
```

---

## Testing Guidelines

### Backend Tests (Go)

```go
func TestInventoryService_CreateResource(t *testing.T) {
    // Arrange
    repo := &mockResourceRepository{}
    service := NewInventoryService(repo, nil, nil)
    ctx := context.Background()

    req := CreateResourceRequest{
        Name:     "Test Resource",
        Quantity: 100,
    }

    // Act
    resource, err := service.CreateResource(ctx, req)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, resource)
    assert.Equal(t, "Test Resource", resource.Name)
    assert.Equal(t, 100, resource.Quantity)
}
```

### Frontend Tests (React)

```typescript
import { render, screen, fireEvent } from '@testing-library/react';
import { ResourceCard } from './ResourceCard';

describe('ResourceCard', () => {
  it('renders resource details', () => {
    const resource = {
      id: '123',
      name: 'AK-74M',
      quantity: 45,
    };

    render(<ResourceCard resource={resource} onEdit={() => {}} onDelete={() => {}} />);

    expect(screen.getByText('AK-74M')).toBeInTheDocument();
    expect(screen.getByText('Quantity: 45')).toBeInTheDocument();
  });

  it('calls onEdit when edit button clicked', () => {
    const onEdit = jest.fn();
    const resource = { id: '123', name: 'Test', quantity: 10 };

    render(<ResourceCard resource={resource} onEdit={onEdit} onDelete={() => {}} />);

    fireEvent.click(screen.getByText('Edit'));

    expect(onEdit).toHaveBeenCalledWith('123');
  });
});
```

---

## Questions?

- 📧 Email: markostrutinsky@example.com
- 💬 GitHub Discussions: [Discussions](https://github.com/markostrutinsky/diploma/discussions)
- 🐛 Issues: [Issues](https://github.com/markostrutinsky/diploma/issues)

---

**Дякуємо за ваш внесок! 🙏**

Made with ❤️ for Ukraine 🇺🇦
