# FRESHNESS-11A — proto-common contract checkpoint

## Checkpoint 0 — исходное состояние и write-ahead решение

- Время старта: `2026-08-15 22:53:38 MSK`.
- Репозиторий: `proto-common`.
- Ветка: `master`.
- Starting HEAD: `de77a16dbff312c16156971d292f7d708b2a197c`
  (`v0.26.0`, `feat(events): add recommendation shelf identity`).
- Исходный `git status --short`: clean.
- Write-scope этого checkpoint: только `events/events.proto`, сгенерированный
  `events/events.pb.go` и этот changelog.
- Tag, push и изменения потребителей запрещены до проверки координатором.

### Ownership и граница события

- `recommendation-service` — единственный authoritative producer tier,
  selection lifecycle и причин отбора.
- `search-service` получает только атомарную read-проекцию и не вычисляет tier.
- Событие не содержит catalog metadata, vectors, vibe tags, summaries или
  freshness source payloads.
- Одна запись события полностью заменяет пару `tier + selection_status` и
  связанные projection fields; частичные обновления запрещены.
- `state_version` монотонен для одного `work_id`. Consumer применяет только
  версию строго больше сохранённой; равная версия идемпотентна, меньшая stale.
- Tombstone использует тот же ordering: `removed=true` удаляет discovery
  projection только если его `state_version` новее сохранённого.

### Точное proto-решение до изменения кода

Новые enum и wire values:

```text
WorkDiscoveryTier
  WORK_DISCOVERY_TIER_UNSPECIFIED = 0
  WORK_DISCOVERY_TIER_COLD = 1        // domain numeric tier 0
  WORK_DISCOVERY_TIER_DISCOVERY = 2   // domain numeric tier 1
  WORK_DISCOVERY_TIER_READABLE = 3    // domain numeric tier 2

WorkDiscoverySelectionStatus
  WORK_DISCOVERY_SELECTION_STATUS_UNSPECIFIED = 0
  WORK_DISCOVERY_SELECTION_STATUS_NOT_APPLICABLE = 1
  WORK_DISCOVERY_SELECTION_STATUS_ACTIVE = 2
  WORK_DISCOVERY_SELECTION_STATUS_GRACE = 3
  WORK_DISCOVERY_SELECTION_STATUS_INACTIVE = 4
  WORK_DISCOVERY_SELECTION_STATUS_BLOCKED = 5

WorkDiscoveryReason
  WORK_DISCOVERY_REASON_UNSPECIFIED = 0
  WORK_DISCOVERY_REASON_OL_TRENDING = 1
  WORK_DISCOVERY_REASON_COMMERCE_HISTORY = 2
  WORK_DISCOVERY_REASON_RECENT_RELEASE_ACTIVITY = 3
  WORK_DISCOVERY_REASON_READING_LOG_RECENT = 4
  WORK_DISCOVERY_REASON_ALL_TIME_ACTIVE_CLASSIC = 5
  WORK_DISCOVERY_REASON_LANGUAGE_GENRE_QUOTA = 6
  WORK_DISCOVERY_REASON_READABLE_CONTENT = 7
  WORK_DISCOVERY_REASON_MANUAL_EDITORIAL = 8
  WORK_DISCOVERY_REASON_GRACE_RETENTION = 9
  WORK_DISCOVERY_REASON_OPERATOR_BLOCK = 10

WorkDiscoveryRankBucket
  WORK_DISCOVERY_RANK_BUCKET_UNSPECIFIED = 0
  WORK_DISCOVERY_RANK_BUCKET_NOT_APPLICABLE = 1
  WORK_DISCOVERY_RANK_BUCKET_TOP = 2
  WORK_DISCOVERY_RANK_BUCKET_HIGH = 3
  WORK_DISCOVERY_RANK_BUCKET_MEDIUM = 4
  WORK_DISCOVERY_RANK_BUCKET_BASE = 5
```

Raw discovery score не публикуется: search получает устойчивый controlled
rank bucket, а его границы принадлежат `policy_version`.

Поля `WorkDiscoveryUpdatedEvent` и номера:

```text
1  string event_id
2  string work_id
3  WorkDiscoveryTier tier
4  WorkDiscoverySelectionStatus selection_status
5  repeated WorkDiscoveryReason reasons
6  bool discovery_eligible
7  WorkDiscoveryRankBucket rank_bucket
8  string policy_version
9  uint64 state_version
10 google.protobuf.Timestamp occurred_at
11 bool removed
```

Validation contract для будущих producer/consumer:

- при `removed=false`: IDs, `policy_version`, `state_version`, `occurred_at`,
  tier и selection status обязательны; `UNSPECIFIED` недопустим;
- при `removed=true`: consumer удаляет projection; state payload игнорируется
  и producer должен оставлять tier/status/reasons/eligibility/rank значениями
  по умолчанию;
- `NOT_APPLICABLE` явно кодирует readable path без selection lifecycle;
- rank bucket не является tier и не влияет на eligibility самостоятельно.

### План проверок

1. Изменить только `events/events.proto`.
2. Выполнить `make gen` и проверить generated diff.
3. Выполнить `go test ./...` и `go vet ./...`.
4. Проверить protobuf descriptor/field numbers и отсутствие неожиданных файлов.
5. Проверить `git diff --check`, `git status`, отсутствие `replace` и бинарных
   build-артефактов.
6. Записать фактические результаты в этот changelog, сделать один локальный
   commit и остановиться без tag/push.

## Checkpoint 1 — proto contract записан, перед генерацией

- В `events/events.proto` добавлены ровно четыре заранее зафиксированных enum
  и `WorkDiscoveryUpdatedEvent` с номерами полей 1–11 без отклонений от
  Checkpoint 0.
- Событие размещено в search-event разделе, но authoritative producer явно
  указан как recommendation-service.
- Поля catalog metadata, vectors, generated signals и raw source observations
  не добавлялись.
- Следующий существенный шаг: `make gen`, затем отдельная проверка, что
  generated diff содержит только ожидаемые enum/message descriptors и Go API.

## Checkpoint 2 — первый `make gen` не стартовал

- Команда: `make gen`.
- Результат: fail до изменения generated files — `protoc-gen-go: program not
  found or is not executable`.
- Причина: `/Users/kirimatt/go/bin/protoc-gen-go` существует, но GOPATH bin не
  входит в текущий `PATH`.
- Repository diff после сбоя содержит только changelog и исходный proto.
- Следующий шаг: повторить неизменённый `make gen` с одноразовым
  `PATH="$(go env GOPATH)/bin:$PATH"`; копировать или собирать бинарник в
  репозиторий не требуется.

## Checkpoint 3 — генерация выполнена, generated diff требует аудита

- Команда: `PATH="$(go env GOPATH)/bin:$PATH" make gen`.
- Результат: success; изменились `events/events.proto` и
  `events/events.pb.go`, `common/common.pb.go` не изменился.
- Наблюдение: generated diff составляет `1136` строк (`890 insertions`, `316
  deletions` суммарно с proto), что существенно больше ожидаемого локального
  добавления.
- До проверок/commit необходимо сравнить generator versions и отделить
  descriptor/API additions от несвязанного formatter/toolchain churn. Большой
  generated diff без объяснения не принимается.

## Checkpoint 4 — причина churn найдена, контракт переносится в конец файла

- `protoc-gen-go v1.36.1` и generated header совпадают с исходным файлом;
  toolchain drift отсутствует.
- Большой diff вызван размещением новых enum/message перед существующими
  catalog definitions: generator перенумеровал внутренние enum/message indexes,
  dependency indexes и exporter tables. Wire field numbers старых сообщений
  не менялись, но review surface стал неоправданно большим.
- Перед следующей правкой решение: переместить неизменённый новый контракт в
  отдельный `Recommendation -> Search Discovery Projection Events` раздел в
  конец `events.proto`. Это не меняет ни одного зафиксированного enum value или
  field number и должно сохранить индексы всех существующих generated types.
- После перемещения `make gen` выполняется повторно, а diff проверяется заново.

## Checkpoint 5 — generated diff принят к финальной проверке

- Повторный `make gen` успешен.
- Existing message order сохранён; новый message добавлен последним. Existing
  enum order также сохранён, новые enum добавлены после него.
- Оставшийся размер generated diff (`695 insertions`, `191 deletions` только в
  `events.pb.go`) объясняется ожидаемыми четырьмя enum, одним message/getters,
  protobuf raw descriptor replacement, увеличением enum/message tables и
  сдвигом Go dependency indexes из-за добавления enum types. Public definitions
  существующих сообщений в diff не меняются.
- Следующий существенный шаг: compile/static checks, descriptor/field audit,
  whitespace/status/artifact audit. После них changelog получит фактические
  результаты до commit.

## Checkpoint 6 — финальная проверка перед commit

Фактически выполнено:

- `go test ./...` — PASS; все три packages компилируются, test files нет.
- `go vet ./...` — PASS.
- `protoc -I . --descriptor_set_out=/tmp/proto-common-events.pb
  --include_imports events/events.proto` — PASS.
- повторный `PATH="$(go env GOPATH)/bin:$PATH" make gen` — PASS и не изменил
  SHA-256 generated файла, генерация детерминирована.
- `git diff --check` — PASS.
- Ручная сверка proto и generated struct tags — поля 1–11 и типы совпадают с
  Checkpoint 0.
- Existing generated public message definitions не менялись; additions и
  descriptor/index tables соответствуют четырём enum и одному event.
- `rg '^replace ' go.mod` — пусто; local proto replace отсутствует.
- executable-file scan вне `.git` — пусто; бинарников/build-артефактов в
  репозитории нет.
- Final pre-commit status содержит только ожидаемые:
  `events/events.proto`, `events/events.pb.go`,
  `FRESHNESS-11A-CHANGELOG.md`.

Контрольные SHA-256 перед commit:

```text
c53d006078b3d91a8defbdec8c6c03da932389118ea84668cf5325320ed8b7b9  events/events.proto
c4d50b0a3f8b410c0f59cfffd786c168e7eaa666dc1ca4b8c963ce2ec513f88a  events/events.pb.go
```

Следующий и последний шаг checkpoint: один локальный commit. Tag и push не
выполняются; дальнейшая реализация 11A остаётся координатору и другим
репозиториям после contract review.
