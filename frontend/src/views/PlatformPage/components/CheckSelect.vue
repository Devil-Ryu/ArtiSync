<template>
    <div>
        <t-input-adornment>
            <template #prepend>
                <t-button variant="outline" @click="props.item.Disabled = !props.item.Disabled">
                    <div v-if="!props.item.Disabled" style="color:cornflowerblue;">T</div>
                    <div v-if="props.item.Disabled" style="color:chocolate;">F</div>
                </t-button>
            </template>
            <t-select class="c-select"  :options="['Headers', 'Query'].map((value)=>{return {value, label:value}})" v-model:value="props.item.Type"/>
        </t-input-adornment>
    </div>
</template>

<script setup>
const props = defineProps({
    item: {
        type: Object,
        default: () => {
            return {
                Disabled: {
                    type: Boolean,
                    default: false,
                },
                Type: {
                    type: String,
                    default: "",
                }
            }
        }
    }

})
</script>

<style>
.c-select {
    width: 100px;
}
</style>