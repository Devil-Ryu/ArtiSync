<template>
    <div ref="messageRef" id="messageArea"></div>
</template>

<script setup>
import { EditorState } from "@codemirror/state"
import { foldGutter, syntaxHighlighting, foldAll, defaultHighlightStyle } from "@codemirror/language"
import { json } from "@codemirror/lang-json"
import { EditorView, highlightSpecialChars, drawSelection, lineNumbers } from "@codemirror/view"
import { onMounted, ref, watch } from "vue";

const props = defineProps({
    content: {
        type: String,
        default: `{"fileName":"AGREEMENT.md", "filePath": "/Users/ryu/MyProjects/Desktop_Applications/WailsProjects/ArtiSync", "imagePath":"/Users/ryu/MyProjects/Desktop_Applications/WailsProjects/ArtiSync", "imageNum":{"test":"dddd"}}`,
    },
    contentType: {
        type: String,
        default: "UTF-8"
    }
})

const messageRef = ref()
const view = ref()



let myTheme = EditorView.theme({
    "&": { height: "300px" },
    ".cm-scroller": { overflow: "auto" }
})


// 加载时挂载，进行初始化
onMounted(() => {
    let startState = EditorState.create({
        extensions: [
            // syntaxHighlighting(defaultHighlightStyle),
            // EditorView.editable.of(false),  // 不可编辑
            // EditorView.lineWrapping,
            // json(),
            myTheme,
        ],
        doc: props.content,
    })
    view.value = new EditorView({
        state: startState,
        parent: messageRef.value

    })
})

watch(
    () => props,
    (newValue, oldValue) => {
        if (view.value !== undefined) {
            view.value.destroy()
            console.log("destroy")
        }
        // nextTick(() => {
        //     console.log("messageRef.value.scrollTop1", messageRef.value.scrollTop)
        //     // messageRef.value.scrollTo({top:0})
        // })

        // messageRef.value.scrollTop = 0
        let startState = EditorState.create({
            extensions: [
                // syntaxHighlighting(defaultHighlightStyle),
                // EditorView.editable.of(false),  // 不可编辑
                // EditorView.lineWrapping,
                // json(),
                myTheme,
            ],
            doc: props.content,
        })

        view.value = new EditorView({
            state: startState,
            parent: messageRef.value

        })

    },
    { deep: true }
)

</script>