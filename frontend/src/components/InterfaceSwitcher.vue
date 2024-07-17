<script lang="ts" setup>
import TitleSetter from './TitleSetter.vue'
import {DebugLog, Reset, Title} from "../../wailsjs/go/main/App";
import WordProcessor from "./WordProcessor.vue";
import { titleData } from "./TitleSetter.vue"
import {reactive} from "vue";


let novel = reactive({
  Title: "untitled",
  titleSet: false,
});

function exportDebugLog() {
  DebugLog("", true)
}

function getTitle(){
  Title().then(result => {
    novel.Title = result
    if (novel.Title != "" && novel.Title != "untitled" && novel.Title != "TBA") {
      novel.titleSet = true
    }
  })
}

function resetNovel(){
  novel.Title = "untitled"
  novel.titleSet = false
  titleData.ready = false
  Reset()
}

</script>

<template>
  <TitleSetter v-if="!titleData.ready" />
  <button class="btn" @click="resetNovel" v-if="titleData.ready">Delete it all</button>
  <WordProcessor v-if="titleData.ready" />
  <button class="btn" @click="exportDebugLog">Export Debug Log</button>
</template>

<style scoped>

</style>