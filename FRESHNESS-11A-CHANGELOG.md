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

## Checkpoint 7 — возврат контракта после cross-repo аудита

- Время повторного старта: `2026-08-15`, после локального commit `8d15f6d`.
- Starting HEAD этого доработочного checkpoint:
  `8d15f6dde666d5056e1cac8395794c822f26d1f2`.
- Starting `git status --short`: clean.
- Write-scope остаётся только `proto-common`: этот changelog,
  `events/events.proto`, generated `events/events.pb.go`.
- Tag/push и изменения producer/consumer репозиториев запрещены.

### Результат read-only cross-repo аудита

Проверены текущие локальные состояния:

- `search-service` HEAD `e9c28175cc99cfe834e95d951b8a9a115c1d6e1a`:
  audit подтверждает, что будущему sole online writer не хватает
  `title_normalized`, `author_normalized`, полного ISBN set, cover, language и
  монотонной catalog facts version. Pending metadata/discovery halves уже
  спроектированы независимо, но безопасного restore/unpublish correlation в
  первом draft события нет.
- `catalog-service` HEAD `5559b06186d0c7548344cdc691c7c4bd23506aad`:
  `WorkUpsertedEvent` публикуется transactional River outbox, но job/event
  содержит только поля 1–13. Projection replay пересобирает snapshot и ставит
  новый `time.Now()`, поэтому `occurred_at` не является authoritative ordering.
  `WorkPurgedEvent` также не несёт ordering version. Catalog model уже содержит
  normalized title/author, cover/language, а canonical ISBN set можно собрать
  из editions; дублировать typed labels отдельным proto type не требуется.
- `recommendation-service` HEAD
  `1174224ddc809045db09325bbc2eb7929c5fdf15`: catalog projection сейчас
  отбрасывает stale snapshot по `occurred_at` и не может записать/передать
  source catalog version в discovery state. Поэтому discovery, рассчитанный по
  metadata до unpublish/restore, способен снова стать видимым после прихода
  другой половины проекции.

Вывод: одного `WorkDiscoveryUpdatedEvent.state_version` недостаточно. Нужны
независимый strict ordering catalog half и явная версия catalog snapshot,
которую использовал reco. Search показывает discovery half только при точном
равенстве source и current metadata versions.

### Точный planned wire diff до proto-правок

`WorkUpsertedEvent` сохраняет поля 1–13 и получает:

```text
14 uint64 metadata_version
15 string title_normalized
16 string author_normalized
17 repeated string isbn13s
18 string cover_url
19 string language
```

Semantics:

- `metadata_version` строго монотонна per `work_id` и версионирует полный
  catalog metadata snapshot;
- incoming greater version атомарно заменяет только всю metadata half;
  equal — idempotent, lower — stale;
- `isbn13s` — полный canonical sorted/deduplicated ISBN-13 set всех editions;
- `labels` остаются display names; typed facts остаются полным JSON projection
  в `metadata["label_projection_v1"]`;
- metadata event не изменяет discovery half.

`WorkPurgedEvent` сохраняет поля 1–4 и получает:

```text
5 uint64 metadata_version
```

Это ordered hard tombstone: consumer применяет только более высокую version и
сохраняет tombstone ordering watermark, чтобы старый `WorkUpsertedEvent` не
воскресил удалённый work.

`WorkDiscoveryUpdatedEvent` сохраняет поля 1–11 и получает:

```text
12 uint64 source_metadata_version
```

Semantics:

- версия означает catalog snapshot, на основании которого reco вывел весь
  discovery state;
- greater `state_version` атомарно заменяет только всю discovery half; equal —
  idempotent, lower — stale;
- discovery eligibility видима только при
  `source_metadata_version == metadata_version` актуальной metadata half;
- любое неравенство (source старее или новее) остаётся pending до получения
  совпадающей второй проекции;
- discovery update не изменяет metadata half.

Никаких deprecated/compatibility fields, vectors, generated metadata или raw
discovery score не добавляется. Следующий шаг — только описанный proto diff,
после него отдельный checkpoint перед `make gen`.

## Checkpoint 8 — proto расширен, перед code generation

- `WorkUpsertedEvent` получил только planned fields 14–19; поля 1–13 и их wire
  numbers не менялись.
- Комментарий `labels` уточнён: это display names, а typed facts остаются в
  `metadata["label_projection_v1"]`.
- `WorkPurgedEvent.metadata_version = 5` использует ту же ordering sequence.
- `WorkDiscoveryUpdatedEvent.source_metadata_version = 12` связывает две
  independently versioned projection halves.
- В proto прямо записаны full-half replacement, greater/equal/lower ordering,
  purge watermark и exact-version visibility semantics.
- Отклонений от Checkpoint 7, deprecated fields и compatibility additions нет.

Следующий шаг: `make gen` с уже установленным `protoc-gen-go v1.36.1`, затем
аудит generated public structs/descriptors до запуска compile checks.

## Checkpoint 9 — generated API audit и tightening version invariants

- `make gen` — success на `protoc-gen-go v1.36.1`; изменился только ожидаемый
  `events/events.pb.go`.
- Generated structs/getters подтверждают wire fields:
  `WorkUpsertedEvent` 14–19, `WorkPurgedEvent` 5,
  `WorkDiscoveryUpdatedEvent` 12.
- Большой raw descriptor hunk ожидаем: комментарии и новые поля меняют encoded
  source descriptor; existing field numbers не сдвинулись.
- Перед финальной генерацией требуется comment-only tightening без wire diff:
  все ordering versions (`metadata_version`, `state_version`,
  `source_metadata_version`) должны быть `> 0`; producer не имеет права
  выпускать разные snapshots с одинаковой version, а consumer может считать
  такой конфликт invalid вместо молчаливой перезаписи.

После comment-only правки выполняется финальный `make gen` и полный validation
gate; номера/типы полей не меняются.

## Checkpoint 10 — доработка готова к отдельному локальному commit

Финальный validation gate:

- `PATH="$(go env GOPATH)/bin:$PATH" make gen` — PASS.
- `gofmt -d events/events.pb.go` — пусто; generated Go уже отформатирован.
- `go test ./...` — PASS, три packages компилируются.
- `go vet ./...` — PASS.
- descriptor build в `/tmp/proto-common-events-v027.pb` — PASS.
- повторный `make gen` не изменил SHA-256 `events/events.pb.go` — generated
  output детерминирован.
- ручная сверка source/generated подтвердила planned field numbers и Go types.
- `git diff --check` — PASS.
- `go.mod` не содержит `replace`.
- executable-file scan вне `.git` пуст; repository binaries отсутствуют.
- Final diff относительно `8d15f6d` ограничен тремя ожидаемыми файлами:
  changelog, `events/events.proto`, `events/events.pb.go`.

Контрольные SHA-256 перед commit:

```text
e2d80ccdc838e5f42c326222aa3fe026cdbd6078b9b37b8eafae4c8eb198e7f2  events/events.proto
0e9036045b5b8402a6de57d31aef1d80f8dc3886bcdaede4bad56ff22ce48724  events/events.pb.go
76690c8757d00e153b58fbef668cb97ce04f9dcf820289c1c07f802e66192339  descriptor set (temporary /tmp evidence)
```

Commit intent: отдельный локальный commit поверх `8d15f6d`, без tag/push.
После commit контракт останавливается на coordinator review перед `v0.27.0`.
