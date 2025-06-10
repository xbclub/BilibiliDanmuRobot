<script setup lang="ts">
// import { useI18n } from "vue-i18n";
import Motion from "../utils/motion";
import ReQrcode from "@/components/ReQrcode";
import { useUserStoreHook } from "@/store/modules/user";
import {onMounted, reactive} from "vue";
import {Userlogin, Getlogin, GetloginStatus} from "../../../../wailsjs/go/main/App"
import {bool} from "vue-types";
import router from "@/router";
import {getTopMenu} from "@/router/utils";
// const { t } = useI18n();
const data = reactive({
  qrcode: "",
  islogin: false,
  loginfailed: false
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
    switch (res){
      case 1:
        data.islogin = true
        window.localStorage.setItem("userInfo","true")
        router.push("/welcome")
        return
      case 3:
        getqrcode()
    }
  })
}
async function getqrcode() {
  await GetloginStatus().then(res => {
    data.islogin = res
  })
  if (data.islogin) {
    window.localStorage.setItem("userInfo","true")
    router.push("/welcome")
  }
  Userlogin().then(result => {
    data.qrcode = result.data.url
  })
}
</script>

<template>
  <Motion class="-mt-2 -mb-2"> <ReQrcode :text=data.qrcode /> </Motion>
  <Motion :delay="200" v-if="data.loginfailed">
    <el-text>
      <p class="text-red-500 text-s"   style="white-space: pre-wrap;">登录失败请重新扫码</p>
    </el-text>
  </Motion>
  <Motion :delay="200">
    <el-divider>
      <p class="text-gray-500 text-xs" style="white-space: pre-wrap;">哔哩哔哩APP<br>扫码登录</p>
    </el-divider>
  </Motion>
      <Motion :delay="150" >
        <el-text>
         <p class="text-red-500 text-m"   style="white-space: pre-wrap;">本项目永久免费!!!!</p>
        </el-text>
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
