<script lang="ts">
import {reactive} from 'vue'

export const titleData = reactive({
  title: "untitled",
  ready: false
})

export default {
  data() {
    return {
      titleData: titleData
  }
  }
}
</script>

<script lang="ts" setup>
import {SetTitle} from "../../wailsjs/go/main/App";

function setTitle() {
  console.log("test")
  SetTitle(titleData.title)
  titleData.ready = true
  console.log("set")
}
</script>

<template>
  <main>
    <div id="result" class="result">
      <p>Welcome to Gwab. Let's get writing! Choose a name for your story:</p></div>
    <div id="input" class="input-box">
      <input id="name" v-model="titleData.title" autocomplete="off" class="input" type="text"/>
      <button class="btn" @click="setTitle">Set</button>
    </div>
  </main>
</template>

<style scoped>
.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
}

.input-box .btn {
  width: 60px;
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  margin: 0 0 0 20px;
  padding: 0 8px;
  cursor: pointer;
}

.input-box .btn:hover {
  background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
  color: #333333;
}

.input-box .input {
  border: none;
  border-radius: 3px;
  outline: none;
  height: 30px;
  line-height: 30px;
  padding: 0 10px;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
}

.input-box .input:hover {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.input-box .input:focus {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}
</style>
