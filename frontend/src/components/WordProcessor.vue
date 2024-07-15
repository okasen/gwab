<script setup lang="ts">

import { Save } from '../../wailsjs/go/main/App'
import Editor from "./Editor.vue";
import {reactive, watch} from "vue";

function saveNovel() {
  Save(JSON.parse(data.content).html)
}

let data = reactive({
    content: "",
    wordcount: 0
}
)

function countWords(text: string) {
  console.log(text)
  text = JSON.parse(text).text
  let splitText = text ? text.split(" ") : []
  splitText = splitText.filter((word: string) => word != "" && !(word.startsWith("<") && word.endsWith(">")))
  console.log(splitText)
  data.wordcount = splitText.length;
}

watch(data, () => countWords(data.content || ""))

</script>

<template>
    <button class="btn" @click="saveNovel">Save</button>
  <div class="word-container">
    <editor v-model="data.content"/>
  </div>
  <p>{{ data.wordcount }}</p>
</template>

<style scoped>
  .word-container {
    padding: 1em;
  }
</style>