<template>
    <div ref="codeViewRef" id="codemirror" ></div>
</template>

<script setup>
import { EditorState } from "@codemirror/state"
import { minimalSetup, basicSetup } from "codemirror"
import { foldGutter, syntaxHighlighting, foldAll,defaultHighlightStyle } from "@codemirror/language"
import { json } from "@codemirror/lang-json"
import { EditorView, highlightSpecialChars, drawSelection, lineNumbers} from "@codemirror/view"
import { onMounted, ref, watch } from "vue";

const props = defineProps({
    content: {
        type: String,
        default: `{"fileName":"AGREEMENT.md", "filePath": "/Users/ryu/MyProjects/Desktop_Applications/WailsProjects/ArtiSync", "imagePath":"/Users/ryu/MyProjects/Desktop_Applications/WailsProjects/ArtiSync", "imageNum":{"test":"dddd"}}`,
    },
    contentType: {
        type: String,
        default:  "TEXT"
    }
})

const codeViewRef = ref()
const view = ref()


watch(
  ()=> props,
  (newValue, oldValue) => {
    if (view.value !== undefined) {
        view.value.destroy()
    }

    var content = props.content
    if (props.contentType === "JSON") {
        content = JSON.stringify(JSON.parse(props.content), null, '\t')
    }
    
    let startState = EditorState.create({
        extensions: [
            // minimalSetup,
            // foldGutter(), // 可折叠
            syntaxHighlighting(defaultHighlightStyle),
            EditorView.editable.of(false),  // 不可编辑
            EditorView.lineWrapping,
            json(),
            
        ],
        doc: content,
    })

    view.value = new EditorView({
        state: startState,
        parent: codeViewRef.value

    })

  },
  {deep: true}
)

</script>