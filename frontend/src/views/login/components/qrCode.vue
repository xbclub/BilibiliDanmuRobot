<script setup lang="ts">
// import { useI18n } from "vue-i18n";
import Motion from "../utils/motion";
import ReQrcode from "@/components/ReQrcode";
import { useUserStoreHook } from "@/store/modules/user";
import {onMounted, reactive} from "vue";
import {Userlogin,Getlogin} from "../../../../wailsjs/go/main/App"
import {bool} from "vue-types";
// const { t } = useI18n();
const data = reactive({
  qrcode: "",
  islogin: false
})
onMounted(()=>{
  getqrcode()
  const id = setInterval(() => {
    checkLogin()
    if (data.islogin) {
      clearInterval(id)
    }
  }, 5000)
})
function checkLogin(){
  Getlogin().then(res => {
    data.islogin = res
  })
}
function getqrcode() {
  Userlogin().then(result => {
    data.qrcode = result.data.url
  })
}
</script>

<template>
  <Motion class="-mt-2 -mb-2"> <ReQrcode :text=data.qrcode /> </Motion>
  <Motion :delay="200">
    <el-divider>
      <p class="text-gray-500 text-xs"  style="white-space: pre-wrap;">哔哩哔哩APP<br>扫码登录</p>
    </el-divider>
  </Motion>
  <Motion :delay="150">
    <el-button
      class="w-full mt-4"
      @click=getqrcode
    >
<!--      {{ t("login.back") }}-->
      点击刷新验证码
    </el-button>
  </Motion>
<!--  <Motion :delay="150">-->
<!--    <el-divider>-->
<!--      <p class="text-gray-500 text-xs"  style="white-space: pre-wrap;">当前状态<br>{{data.islogin}}</p>-->
<!--    </el-divider>-->
<!--  </Motion>-->
</template>
