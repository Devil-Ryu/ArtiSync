<template>
    <div>
        <t-input-adornment v-if="!props.item[props.Dynamic]">
            <template #prepend>
                <t-select :placeholder="props.placeholder" :disabled="$attrs.disabled" :inputProps="{ status: $attrs.status }"
                    class="d-select" @change="(value) => { props.item[props.Dynamic] = Boolean(value) }">
                    <template #valueDisplay="{ value, onClose }">
                        <LettersSIcon style="color: cornflowerblue;" />
                    </template>
                    <t-option v-for="item in options" :key="item.value" :value="item.value" :label="item.label">
                        <LettersDIcon v-if="item.value" style="color:chocolate;" />
                        <LettersSIcon v-if="!item.value" style="color: cornflowerblue;" />
                        {{ item.label }}
                    </t-option>
                </t-select>
            </template>
            <t-input class="d-value" v-if="!props.item[props.Dynamic]" v-model="props.item[props.Value]"
                :placeholder="props.placeholder" :status="$attrs.status" :disabled="$attrs.disabled" />
        </t-input-adornment>

        <t-input-adornment v-if="props.item[props.Dynamic]">
            <template #prepend>
                <t-select :placeholder="props.placeholder" :disabled="$attrs.disabled" :inputProps="{ status: $attrs.status }"
                    style="background-color: white; width: 40px; "
                    @change="(value) => { props.item[props.Dynamic] = Boolean(value) }">
                    <template #valueDisplay="{ value, onClose }">
                        <LettersDIcon style="color:chocolate;" />

                    </template>
                    <t-option v-for="item in options" :key="item.value" :value="item.value" :label="item.label">
                        <LettersDIcon v-if="item.value" style="color:chocolate; " />
                        <LettersSIcon v-if="!item.value" style="color: cornflowerblue;" />
                        {{ item.label }}
                    </t-option>
                </t-select>
            </template>
            <t-cascader class="d-value" :placeholder="props.placeholder" :inputProps="{ status: $attrs.status }"
                v-if="props.item[props.Dynamic]" v-model="props.item[props.Value]" :options="props.cascaderOptions"
                check-strictly />
        </t-input-adornment>
    </div>
</template>

<script setup>
import { LettersDIcon, LettersSIcon } from 'tdesign-icons-vue-next';
const props = defineProps({
    Dynamic: {
        type: String,
        default: "Dynamic",
    },
    Value: {
        type: String,
        default: "Value",
    },
    item: {
        type: Object,
        default: () => {
            return {
                Dynamic: false,
                Value: "",
            }
        }
    },
    cascaderOptions: {
        type: Array,
        default: [],
    },
    placeholder: {
        type: String,
        default: "请输入值"
    }
})

// const emits = defineEmits(['update:dynamic', 'update:value'])
const options = [
    { label: "动态", value: 1, description: "自定义下拉选项的选择一段" },
    { label: "静态", value: 0, description: "自定义下拉选项的选择一段" },
]
</script>

<style>
.d-select {
    width: 40px;
    background-color: white;
}

.d-value {
    width: 100%;
}
</style>