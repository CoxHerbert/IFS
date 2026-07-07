<template>
  <div class="app-container build-page">
    <section class="hero-panel">
      <div class="hero-copy">
        <span class="hero-kicker">Form Builder</span>
        <h1>表单构建工作台</h1>
        <p>
          这里用于快速搭建表单结构，支持字段配置、实时预览、JSON 导出以及本地草稿保存，适合作为后续接入后端配置前的前端工作台。
        </p>
        <div class="hero-actions">
          <a-button type="primary" @click="addField('input')">新增单行输入</a-button>
          <a-button @click="addField('select')">新增下拉选择</a-button>
          <a-button @click="saveDraft">保存草稿</a-button>
          <a-button @click="loadDraft">加载草稿</a-button>
        </div>
      </div>

      <div class="hero-side">
        <div class="status-card">
          <div class="status-head">
            <span>当前表单</span>
            <strong>{{ schema.fields.length }} 个字段</strong>
          </div>
          <ul class="status-list">
            <li>支持输入框、文本域、数字框、下拉、单选、多选、日期和开关。</li>
            <li>字段列表可调整顺序、复制、删除，并同步更新右侧属性面板。</li>
            <li>可导出 JSON 结构，便于后续接入接口或动态渲染引擎。</li>
          </ul>
        </div>
      </div>
    </section>

    <a-row :gutter="16" class="workspace-row">
      <a-col :xs="24" :xl="6">
        <section class="panel-card">
          <div class="panel-head">
            <div>
              <h2>基础配置</h2>
              <p>维护表单标题、描述和字段清单。</p>
            </div>
          </div>

          <a-form layout="vertical" :model="schema" class="meta-form">
            <a-form-item label="表单标题">
              <a-input v-model:value="schema.title" placeholder="请输入表单标题" />
            </a-form-item>
            <a-form-item label="表单描述">
              <a-textarea v-model:value="schema.description" :rows="3" placeholder="请输入表单描述" />
            </a-form-item>
          </a-form>

          <div class="action-grid">
            <a-button v-for="item in fieldTypes" :key="item.type" block @click="addField(item.type)">
              {{ item.label }}
            </a-button>
          </div>

          <div class="field-list">
            <div class="list-head">
              <strong>字段列表</strong>
              <a-button type="link" size="small" danger @click="resetSchema">重置</a-button>
            </div>

            <div
              v-for="(field, index) in schema.fields"
              :key="field.id"
              class="field-item"
              :class="{ active: field.id === activeFieldId }"
              @click="selectField(field.id)"
            >
              <div class="field-top">
                <span>{{ index + 1 }}. {{ field.label || "未命名字段" }}</span>
                <a-tag>{{ getFieldTypeLabel(field.type) }}</a-tag>
              </div>
              <p>{{ field.name }}</p>
              <div class="field-actions">
                <a-button size="small" @click.stop="moveField(index, -1)">上移</a-button>
                <a-button size="small" @click.stop="moveField(index, 1)">下移</a-button>
                <a-button size="small" @click.stop="duplicateField(field.id)">复制</a-button>
                <a-button size="small" danger @click.stop="removeField(field.id)">删除</a-button>
              </div>
            </div>

            <a-empty v-if="!schema.fields.length" description="暂无字段，请先添加" />
          </div>
        </section>
      </a-col>

      <a-col :xs="24" :xl="10">
        <section class="panel-card preview-panel">
          <div class="panel-head">
            <div>
              <h2>实时预览</h2>
              <p>按当前配置渲染表单，便于确认结构和展示效果。</p>
            </div>
            <a-space>
              <a-button @click="fillDemoData">填充示例</a-button>
              <a-button type="primary" @click="submitPreview">模拟提交</a-button>
            </a-space>
          </div>

          <div class="preview-shell">
            <div class="preview-header">
              <h3>{{ schema.title }}</h3>
              <p>{{ schema.description || "暂无描述" }}</p>
            </div>

            <a-form
              ref="previewFormRef"
              layout="vertical"
              :model="previewValues"
              class="preview-form"
            >
              <a-row :gutter="16">
                <a-col
                  v-for="field in schema.fields"
                  :key="field.id"
                  :xs="24"
                  :md="field.span || 24"
                >
                  <a-form-item :label="field.label" :required="field.required">
                    <template v-if="field.type === 'input'">
                      <a-input v-model:value="previewValues[field.name]" :placeholder="field.placeholder" />
                    </template>

                    <template v-else-if="field.type === 'textarea'">
                      <a-textarea
                        v-model:value="previewValues[field.name]"
                        :placeholder="field.placeholder"
                        :rows="field.rows || 4"
                      />
                    </template>

                    <template v-else-if="field.type === 'number'">
                      <a-input-number
                        v-model:value="previewValues[field.name]"
                        :placeholder="field.placeholder"
                        style="width: 100%"
                      />
                    </template>

                    <template v-else-if="field.type === 'select'">
                      <a-select
                        v-model:value="previewValues[field.name]"
                        :placeholder="field.placeholder"
                        :options="field.options"
                      />
                    </template>

                    <template v-else-if="field.type === 'radio'">
                      <a-radio-group v-model:value="previewValues[field.name]">
                        <a-radio v-for="option in field.options" :key="option.value" :value="option.value">
                          {{ option.label }}
                        </a-radio>
                      </a-radio-group>
                    </template>

                    <template v-else-if="field.type === 'checkbox'">
                      <a-checkbox-group v-model:value="previewValues[field.name]" :options="field.options" />
                    </template>

                    <template v-else-if="field.type === 'date'">
                      <a-date-picker
                        v-model:value="previewValues[field.name]"
                        style="width: 100%"
                        :placeholder="field.placeholder"
                      />
                    </template>

                    <template v-else-if="field.type === 'switch'">
                      <a-switch v-model:checked="previewValues[field.name]" />
                    </template>
                  </a-form-item>
                </a-col>
              </a-row>

              <a-form-item class="preview-footer">
                <a-space>
                  <a-button type="primary" @click="submitPreview">提交</a-button>
                  <a-button @click="resetPreview">重置</a-button>
                </a-space>
              </a-form-item>
            </a-form>
          </div>
        </section>
      </a-col>

      <a-col :xs="24" :xl="8">
        <section class="panel-card">
          <div class="panel-head">
            <div>
              <h2>属性配置</h2>
              <p>选中字段后，在这里修改展示与行为配置。</p>
            </div>
            <a-button type="link" @click="copySchema">复制 JSON</a-button>
          </div>

          <a-empty v-if="!activeField" description="请选择左侧字段进行配置" />

          <template v-else>
            <a-form layout="vertical" :model="activeField">
              <a-form-item label="字段标题">
                <a-input v-model:value="activeField.label" placeholder="请输入字段标题" />
              </a-form-item>
              <a-form-item label="字段名">
                <a-input v-model:value="activeField.name" placeholder="请输入字段名" />
              </a-form-item>
              <a-form-item label="占位提示">
                <a-input v-model:value="activeField.placeholder" placeholder="请输入占位提示" />
              </a-form-item>
              <a-form-item label="栅格宽度">
                <a-select v-model:value="activeField.span" :options="spanOptions" />
              </a-form-item>
              <a-form-item v-if="activeField.type === 'textarea'" label="文本域行数">
                <a-input-number v-model:value="activeField.rows" :min="2" :max="10" style="width: 100%" />
              </a-form-item>
              <a-form-item label="必填">
                <a-switch v-model:checked="activeField.required" />
              </a-form-item>
            </a-form>

            <template v-if="hasOptions(activeField.type)">
              <div class="options-head">
                <strong>选项配置</strong>
                <a-button size="small" @click="addOption(activeField)">新增选项</a-button>
              </div>
              <div
                v-for="(option, optionIndex) in activeField.options"
                :key="`${activeField.id}-${optionIndex}`"
                class="option-item"
              >
                <a-input v-model:value="option.label" placeholder="标签" />
                <a-input v-model:value="option.value" placeholder="值" />
                <a-button danger @click="removeOption(activeField, optionIndex)">删除</a-button>
              </div>
            </template>
          </template>

          <div class="schema-box">
            <div class="schema-head">
              <strong>Schema 预览</strong>
              <a-space>
                <a-button size="small" @click="formatSchema">格式化</a-button>
                <a-button size="small" @click="copySchema">复制</a-button>
              </a-space>
            </div>
            <pre>{{ schemaText }}</pre>
          </div>
        </section>
      </a-col>
    </a-row>
  </div>
</template>

<script setup name="Build">
import { computed, reactive, ref, watch } from "vue";
import { message, Modal } from "ant-design-vue";

const DRAFT_KEY = "ifs-form-builder-draft";
const fieldSeed = ref(4);
const previewFormRef = ref();

const fieldTypes = [
  { type: "input", label: "单行输入" },
  { type: "textarea", label: "多行文本" },
  { type: "number", label: "数字输入" },
  { type: "select", label: "下拉选择" },
  { type: "radio", label: "单项选择" },
  { type: "checkbox", label: "多项选择" },
  { type: "date", label: "日期选择" },
  { type: "switch", label: "开关" }
];

const spanOptions = [
  { label: "一整行", value: 24 },
  { label: "半行", value: 12 },
  { label: "三分之一", value: 8 }
];

const schema = reactive(createDefaultSchema());
const activeFieldId = ref(schema.fields[0]?.id || "");
const previewValues = reactive(buildPreviewValues(schema.fields));

const activeField = computed(() => schema.fields.find(item => item.id === activeFieldId.value));
const schemaText = computed(() => JSON.stringify(schema, null, 2));

watch(
  () => schema.fields.map(item => item.name),
  () => syncPreviewValues()
);

watch(
  () => schema.fields.length,
  () => {
    if (!schema.fields.length) {
      activeFieldId.value = "";
      clearPreviewValues();
      return;
    }

    if (!schema.fields.some(item => item.id === activeFieldId.value)) {
      activeFieldId.value = schema.fields[0].id;
    }
  }
);

function createDefaultSchema() {
  return {
    title: "客户信息采集",
    description: "用于演示表单构建配置、实时预览与结构导出。",
    fields: [
      createField("input", {
        label: "客户名称",
        name: "customerName",
        placeholder: "请输入客户名称"
      }),
      createField("select", {
        label: "客户等级",
        name: "customerLevel",
        placeholder: "请选择客户等级"
      }),
      createField("textarea", {
        label: "备注说明",
        name: "remark",
        placeholder: "请输入备注信息"
      })
    ]
  };
}

function createField(type, overrides = {}) {
  const nextId = `field_${fieldSeed.value++}`;
  const defaultOptions = [
    { label: "选项一", value: "option1" },
    { label: "选项二", value: "option2" }
  ];

  const base = {
    id: nextId,
    type,
    label: "未命名字段",
    name: nextId,
    placeholder: "请输入内容",
    required: false,
    span: 24
  };

  if (type === "textarea") {
    base.rows = 4;
  }

  if (type === "switch") {
    base.placeholder = "";
  }

  if (hasOptions(type)) {
    base.options = defaultOptions.map(item => ({ ...item }));
    base.placeholder = type === "select" ? "请选择内容" : "";
  }

  if (type === "date") {
    base.placeholder = "请选择日期";
  }

  return Object.assign(base, overrides);
}

function buildPreviewValues(fields) {
  return fields.reduce((acc, field) => {
    acc[field.name] = getInitialValue(field);
    return acc;
  }, {});
}

function getInitialValue(field) {
  if (field.type === "checkbox") {
    return [];
  }

  if (field.type === "switch") {
    return false;
  }

  return undefined;
}

function hasOptions(type) {
  return ["select", "radio", "checkbox"].includes(type);
}

function getFieldTypeLabel(type) {
  return fieldTypes.find(item => item.type === type)?.label || type;
}

function selectField(id) {
  activeFieldId.value = id;
}

function addField(type) {
  const field = createField(type, {
    label: getFieldTypeLabel(type),
    name: `${type}_${Date.now()}`
  });

  schema.fields.push(field);
  previewValues[field.name] = getInitialValue(field);
  activeFieldId.value = field.id;
  message.success("字段已添加");
}

function removeField(id) {
  const index = schema.fields.findIndex(item => item.id === id);
  if (index === -1) {
    return;
  }

  const [removed] = schema.fields.splice(index, 1);
  delete previewValues[removed.name];
  message.success("字段已删除");
}

function duplicateField(id) {
  const source = schema.fields.find(item => item.id === id);
  if (!source) {
    return;
  }

  const copy = JSON.parse(JSON.stringify(source));
  copy.id = `field_${fieldSeed.value++}`;
  copy.name = `${source.name}_copy`;
  copy.label = `${source.label} 副本`;

  const index = schema.fields.findIndex(item => item.id === id);
  schema.fields.splice(index + 1, 0, copy);
  previewValues[copy.name] = getInitialValue(copy);
  activeFieldId.value = copy.id;
  message.success("字段已复制");
}

function moveField(index, offset) {
  const targetIndex = index + offset;
  if (targetIndex < 0 || targetIndex >= schema.fields.length) {
    return;
  }

  const list = [...schema.fields];
  const [current] = list.splice(index, 1);
  list.splice(targetIndex, 0, current);
  schema.fields.splice(0, schema.fields.length, ...list);
}

function addOption(field) {
  if (!field.options) {
    field.options = [];
  }

  field.options.push({
    label: `选项${field.options.length + 1}`,
    value: `option${field.options.length + 1}`
  });
}

function removeOption(field, optionIndex) {
  field.options.splice(optionIndex, 1);
}

function clearPreviewValues() {
  Object.keys(previewValues).forEach(key => {
    delete previewValues[key];
  });
}

function syncPreviewValues() {
  const validNames = schema.fields.map(item => item.name).filter(Boolean);

  validNames.forEach(name => {
    if (!(name in previewValues)) {
      const field = schema.fields.find(item => item.name === name);
      previewValues[name] = getInitialValue(field);
    }
  });

  Object.keys(previewValues).forEach(key => {
    if (!validNames.includes(key)) {
      delete previewValues[key];
    }
  });
}

function fillDemoData() {
  schema.fields.forEach(field => {
    if (field.type === "input" || field.type === "textarea") {
      previewValues[field.name] = `${field.label}示例`;
    } else if (field.type === "number") {
      previewValues[field.name] = 1;
    } else if (field.type === "select" || field.type === "radio") {
      previewValues[field.name] = field.options?.[0]?.value;
    } else if (field.type === "checkbox") {
      previewValues[field.name] = field.options?.slice(0, 1).map(item => item.value) || [];
    } else if (field.type === "switch") {
      previewValues[field.name] = true;
    } else if (field.type === "date") {
      previewValues[field.name] = undefined;
    }
  });

  message.success("已填充示例数据");
}

function resetPreview() {
  clearPreviewValues();
  Object.assign(previewValues, buildPreviewValues(schema.fields));
}

function submitPreview() {
  const missingField = schema.fields.find(field => {
    if (!field.required) {
      return false;
    }

    const value = previewValues[field.name];
    if (Array.isArray(value)) {
      return !value.length;
    }

    return value === undefined || value === null || value === "";
  });

  if (missingField) {
    message.warning(`请先填写必填字段：${missingField.label}`);
    return;
  }

  Modal.info({
    title: "表单提交结果",
    width: 720,
    content: () => previewSchemaContent(JSON.stringify(previewValues, null, 2))
  });
}

function previewSchemaContent(text) {
  return h("pre", { class: "modal-json" }, text);
}

function formatSchema() {
  message.success("当前 Schema 已按标准 JSON 格式展示");
}

async function copySchema() {
  try {
    await navigator.clipboard.writeText(schemaText.value);
    message.success("Schema 已复制到剪贴板");
  } catch (error) {
    message.error("复制失败，请检查浏览器剪贴板权限");
  }
}

function saveDraft() {
  localStorage.setItem(DRAFT_KEY, schemaText.value);
  message.success("草稿已保存到本地");
}

function loadDraft() {
  const raw = localStorage.getItem(DRAFT_KEY);
  if (!raw) {
    message.warning("本地没有可用草稿");
    return;
  }

  try {
    const parsed = JSON.parse(raw);
    schema.title = parsed.title || "";
    schema.description = parsed.description || "";
    schema.fields.splice(0, schema.fields.length, ...(parsed.fields || []));
    fieldSeed.value = schema.fields.length + 10;
    clearPreviewValues();
    Object.assign(previewValues, buildPreviewValues(schema.fields));
    activeFieldId.value = schema.fields[0]?.id || "";
    message.success("草稿已加载");
  } catch (error) {
    message.error("草稿解析失败");
  }
}

function resetSchema() {
  const next = createDefaultSchema();
  schema.title = next.title;
  schema.description = next.description;
  schema.fields.splice(0, schema.fields.length, ...next.fields);
  activeFieldId.value = schema.fields[0]?.id || "";
  resetPreview();
}
</script>

<style scoped lang="scss">
.build-page {
  min-height: calc(100vh - 84px);
  background:
    radial-gradient(circle at top left, rgba(249, 115, 22, 0.14), transparent 24%),
    radial-gradient(circle at right bottom, rgba(14, 165, 233, 0.12), transparent 26%),
    linear-gradient(180deg, #f8fafc 0%, #eef4f8 100%);
}

.hero-panel {
  display: grid;
  grid-template-columns: minmax(0, 1.8fr) minmax(280px, 0.9fr);
  gap: 16px;
  margin-bottom: 16px;
}

.hero-copy,
.status-card,
.panel-card {
  border: 1px solid rgba(148, 163, 184, 0.2);
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 16px 40px rgba(15, 23, 42, 0.06);
}

.hero-copy {
  padding: 28px 30px;
  border-radius: 24px;
}

.hero-kicker {
  display: inline-flex;
  align-items: center;
  height: 28px;
  padding: 0 12px;
  border-radius: 999px;
  background: rgba(249, 115, 22, 0.12);
  color: #c2410c;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.hero-copy h1 {
  margin: 16px 0 10px;
  color: #0f172a;
  font-size: 34px;
  line-height: 1.15;
}

.hero-copy p {
  margin: 0;
  max-width: 780px;
  color: #475569;
  line-height: 1.8;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 22px;
}

.status-card {
  height: 100%;
  padding: 22px;
  border-radius: 24px;
  background: linear-gradient(180deg, rgba(15, 23, 42, 0.96) 0%, rgba(30, 41, 59, 0.94) 100%);
  color: #e2e8f0;
}

.status-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 18px;
}

.status-head span {
  color: rgba(226, 232, 240, 0.8);
  font-size: 14px;
}

.status-head strong {
  display: inline-flex;
  align-items: center;
  height: 30px;
  padding: 0 12px;
  border-radius: 999px;
  background: rgba(56, 189, 248, 0.18);
  color: #bae6fd;
  font-size: 13px;
}

.status-list {
  margin: 0;
  padding-left: 18px;
  line-height: 1.9;
}

.workspace-row {
  align-items: stretch;
}

.panel-card {
  height: 100%;
  padding: 22px;
  border-radius: 24px;
}

.panel-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}

.panel-head h2 {
  margin: 0 0 6px;
  color: #0f172a;
  font-size: 22px;
}

.panel-head p {
  margin: 0;
  color: #64748b;
  line-height: 1.7;
}

.meta-form {
  margin-bottom: 18px;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  margin-bottom: 18px;
}

.field-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.list-head,
.schema-head,
.options-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.field-item {
  padding: 14px;
  border-radius: 18px;
  border: 1px solid rgba(148, 163, 184, 0.16);
  background: #f8fafc;
  cursor: pointer;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
}

.field-item.active {
  border-color: rgba(249, 115, 22, 0.42);
  box-shadow: 0 14px 28px rgba(249, 115, 22, 0.12);
  transform: translateY(-1px);
}

.field-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
  color: #0f172a;
  font-weight: 600;
}

.field-item p {
  margin: 0 0 12px;
  color: #64748b;
  font-size: 13px;
}

.field-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.preview-panel {
  display: flex;
  flex-direction: column;
}

.preview-shell {
  flex: 1;
  padding: 22px;
  border-radius: 22px;
  background: linear-gradient(180deg, #fff 0%, #f8fafc 100%);
  border: 1px solid rgba(148, 163, 184, 0.16);
}

.preview-header {
  margin-bottom: 16px;
}

.preview-header h3 {
  margin: 0 0 8px;
  color: #0f172a;
  font-size: 24px;
}

.preview-header p {
  margin: 0;
  color: #64748b;
}

.preview-footer {
  margin-bottom: 0;
}

.option-item {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr) auto;
  gap: 8px;
  margin-bottom: 10px;
}

.schema-box {
  margin-top: 20px;
  padding: 16px;
  border-radius: 18px;
  background: #0f172a;
  color: #e2e8f0;
}

.schema-box pre {
  margin: 14px 0 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 12px;
  line-height: 1.7;
}

:deep(.modal-json) {
  max-height: 420px;
  overflow: auto;
  margin: 0;
  padding: 14px;
  border-radius: 12px;
  background: #f8fafc;
}

@media (max-width: 1400px) {
  .hero-panel {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .hero-copy {
    padding: 22px 18px;
  }

  .hero-copy h1 {
    font-size: 28px;
  }

  .panel-card,
  .status-card {
    border-radius: 18px;
  }

  .action-grid,
  .option-item {
    grid-template-columns: 1fr;
  }
}
</style>
