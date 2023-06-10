<script setup lang="ts">
import {onMounted, reactive} from "vue";

defineOptions({
  name: "Welcome"
});
const formRef = ref<FormInstance>()
const data = reactive({
  isrunning: false,
  savestatus: false,
  savemsg:"",
  tabledata: [],
  dialogdata: {
    uid: "",
    msg: ""
  },
  dialogtitle:"",
  dialogVisible: false,
  form: {
    RoomId: 3,
    DanmuLen: 20,
    EntryMsg: "花花机器人进入直播间",
    PKNotice: true,
    InteractWord: true,
    EntryEffect: true,
    ThanksGift: true,
    WelcomeSwitch: true,
    WelcomeString: {
      "123456": "欢迎宇宙无敌最帅的xxx进入直播间"
    },
    WelcomeDanmu: ["欢迎 {user} ~","欢迎 {user} 木嘛~","欢迎 {user} 好诶~"],
    RobotName: "花花",
    TalkRobotCmd: "花花",
    RobotMode: "QingYunKe",
    ChatGPT:{
      APIToken: ""
    },
    FocusDanmu: [
      "啾咪~",
      "喜欢可以领牌牌哦~",
      "么么哒~",
      "入股不亏哦~",
      "贴贴~"
    ],
    CronDanmu: false,
    CronDanmuList: [{
      Cron: "*/2 * * * *",
      Random:true,
      Danmu:[
        "喜欢主播请关注, 主播带你去致富~",
        "万水千山总是情, 上个舰长行不行~",
        "喜欢主播的小伙伴可以动动小手点个关注~",
        "喜欢主播的小伙伴，点点关注不迷路~",
        "你已经是成熟的观众了，该学会自己上船了~",
        "小礼物和弹幕都是对主播的支持哦，比心心~",
        "有一种关心叫关注，有一种惦记叫入粉",
        "有一种陪伴叫: 加入大航海~",
        "iOS端可关注公众号哗哩哗哩直播姬充值~",
        "万水千山总是情，点个关注行不行~"
      ]
    }],
    WelcomeBlacklistWide:[],
    WelcomeBlacklist:[],

  }
})
import { ref } from 'vue'
import type {FormInstance, FormRules, TabsPaneContext} from 'element-plus'
import {ReadConfig, WriteConfig} from "../../../wailsjs/go/main/App";
import {ElNotification} from "element-plus";
import {Monitor, Start, Stop} from "../../../wailsjs/go/main/Program";

const activeName = ref('first')

const handleClick = (tab: TabsPaneContext, event: Event) => {
  console.log(tab, event)
}
const addWelcomeDanmu = () => {
  data.form.WelcomeDanmu.push("{user}")
}
const addWelcomeBlacklist = () => {
  data.form.WelcomeBlacklist.push("")
}
const addWelcomeBlacklistWide = () => {
  data.form.WelcomeBlacklistWide.push("")
}
const deleteWelcomeDanmu = (item: number) => {
  // const index = data.form.WelcomeDanmu.indexOf(item)
  if (item !== -1) {
    data.form.WelcomeDanmu.splice(item, 1)
  }
}
const deleteWelcomeBlacklist = (item: number) => {
  // const index = data.form.WelcomeDanmu.indexOf(item)
  if (item !== -1) {
    data.form.WelcomeBlacklist.splice(item, 1)
  }
}
const deleteWelcomeBlacklistWide = (item: number) => {
  // const index = data.form.WelcomeDanmu.indexOf(item)
  if (item !== -1) {
    data.form.WelcomeBlacklistWide.splice(item, 1)
  }
}
const addFocusDanmu = () => {
  data.form.FocusDanmu.push("")
}
const deleteFocusDanmu = (item: number) => {
  // const index = data.form.WelcomeDanmu.indexOf(item)
  if (item !== -1) {
    data.form.FocusDanmu.splice(item, 1)
  }
}
function formatWelcomeString() {
  const arr = Object.entries(data.form.WelcomeString);
  data.tabledata=arr.map(([key, value]) => ({ key, value }))
}
function addWelcomeString(){
  data.dialogdata={
    uid: "",
    msg: ""
  }
  data.dialogtitle = "新增"
  data.dialogVisible = true

}
function saveWelcomeString(){
  data.form.WelcomeString[data.dialogdata.uid] = data.dialogdata.msg
  formatWelcomeString()
  data.dialogVisible = false
}
function deleteWelcomeString(item:string){
  // console.log(item)
  if (data.form.WelcomeString.hasOwnProperty(item)) {
    delete data.form.WelcomeString[item];
  }
  formatWelcomeString()
}
function addDanmulist(item:string){
  data.form.CronDanmuList.push({
    Cron: "*/2 * * * *",
    Random:true,
    Danmu:[
      "",
    ]
  })
}
function deleteDanmulist(row: number){
  if (row !== -1) {
    data.form.CronDanmuList.splice(row, 1)
  }
}
function addDanmu(item:number){
  data.form.CronDanmuList[item].Danmu.push("")
}
function deleteDanmu(item:number,row: number){
  if (item !== -1&&row!== -1) {
    data.form.CronDanmuList[item].Danmu.splice(row, 1)
  }
}
async function saveConfig() {
  data.savestatus = false
  data.savemsg = ""
  await WriteConfig(JSON.stringify(data.form)).then(res => {
    data.savestatus = res.Code
    data.savemsg = res.Msg
  })
  if (data.savestatus) {
    ElNotification({
      title: '操作成功',
      message: '操作成功，正在重启机器人',
      type: 'success',
    })
    restart()
  }else {
    ElNotification({
      title: '保存失败',
      message: '请修改配置知道保存正确为止，错误信息：'+data.savemsg,
      type: 'error',
    })
  }
}
async function pgstart() {
  await Monitor().then(res => {
    data.isrunning = res
  })
  if (data.isrunning == false) {
    await Start().then(res => {
      data.isrunning = res
    })
  }
}
async function pgstop(){
  await Monitor().then(res => {
    data.isrunning = res
  })
  if (data.isrunning == true) {
    await Stop().then(res => {
      data.isrunning = !res
    })
  }
}
async function restart() {
  await pgstop()
  console.log(data.isrunning)
  if (data.isrunning == false) {
    pgstart()
  }
}
const rules = reactive<FormRules>({
  RoomId:[
    { required: true, message: '此项必填' },
    { type: 'number', message: '此处必须为数字' },
  ],
  DanmuLen:[
    { required: true, message: '此项必填' },
    { type: 'number', message: '此处必须为数字' },
  ]
})
onMounted(()=>{

  ReadConfig().then(res=>{
    if (!res.Code){
      console.log(res)
      ElNotification({
        title: '读取配置文件失败',
        message: res.Msg,
        type: 'error',
      })
    }else {
      data.form = res.Form
      formatWelcomeString()
    }
  })

})
</script>

<template>
  <el-form :model="data.form" label-width="120px" ref="formRef" class="el-form" :rules="rules">
    <el-tabs v-model="activeName" class="demo-tabs" type="border-card" @tab-click="handleClick">
      <el-tab-pane label="基础配置" name="first">
        <el-form-item label="直播间号" prop="RoomId">
          <el-input v-model.number="data.form.RoomId" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="可发送弹幕长度"  prop="DanmuLen">
          <el-input v-model.number="data.form.DanmuLen" />
        </el-form-item>
        <el-form-item label="进入直播间消息"  prop="DanmuLen">
          <el-input v-model.number="data.form.EntryMsg" />
        </el-form-item>
        <el-form-item label="欢迎弹幕">
          <el-switch v-model="data.form.InteractWord" />
        </el-form-item>
        <el-form-item label="特效入场欢迎">
          <el-switch v-model="data.form.EntryEffect" />
        </el-form-item>
        <el-form-item label="礼物感谢">
          <el-switch v-model="data.form.ThanksGift" />
        </el-form-item>
        <el-form-item label="PK提醒">
          <el-switch v-model="data.form.PKNotice" />
        </el-form-item>

      </el-tab-pane>

      <el-tab-pane label="欢迎弹幕自定义" name="second">
        <el-form-item>
          <el-tag>Tips: 自定义欢迎语 {user}为用户昵称占位符，列表为随机列表</el-tag>
        </el-form-item>
        <el-form-item>
        <el-button type="primary" @click="addWelcomeDanmu">新增</el-button>
        </el-form-item>
        <el-form-item
          v-for="(items, index) in data.form.WelcomeDanmu"
          :label="'欢迎弹幕 ' + (index+1)"
        >
          <el-input v-model="data.form.WelcomeDanmu[index]" />
          <el-button class="mt-2"
                     @click="deleteWelcomeDanmu(index)"
          >删除</el-button
          >
        </el-form-item>
      </el-tab-pane>
      <el-tab-pane label="指定人欢迎" name="third">
<!--        <template>-->
        <el-form-item label="功能开关">
         <el-switch v-model="data.form.WelcomeSwitch" />
        </el-form-item>
          <el-button @click="addWelcomeString" type="primary">添加</el-button>
          <el-table :data="data.tabledata" border>
            <el-table-column prop="key" label="用户uid"  />
            <el-table-column prop="value" label="欢迎语"  />
            <el-table-column label="操作">
              <template #default="scope">
                  <el-button>编辑</el-button>
                  <el-button @click="deleteWelcomeString(scope.row.key)">删除</el-button>
              </template>
            </el-table-column>
<!--            <el-table-column prop="address" label="Address" />-->
          </el-table>
<!--        </template>-->
      </el-tab-pane>
      <el-tab-pane label="ai机器人" name="fourth">
        <el-form-item label="机器人名称">
          <el-input v-model.number="data.form.RobotName" />
        </el-form-item>
        <el-form-item label="触发关键字">
          <el-input v-model.number="data.form.TalkRobotCmd" />
        </el-form-item>
        <el-form-item label="服务提供商" prop="region">
          <el-select v-model="data.form.RobotMode" placeholder="请选择">
            <el-option label="QingYunKe" value="QingYunKe" />
            <el-option label="ChatGPT" value="ChatGPT" />
          </el-select>
        </el-form-item>
      </el-tab-pane>
      <el-tab-pane label="关注答谢语" name="fifth">
        <el-form-item>
          <el-button type="primary" @click="addFocusDanmu">新增</el-button>
        </el-form-item>
        <el-form-item
          v-for="(items, index) in data.form.FocusDanmu"
          :label="'关注答谢语 ' + (index+1)"
        >
          <el-input v-model="data.form.FocusDanmu[index]" />
          <el-button class="mt-2"
                     @click="deleteFocusDanmu(index)"
          >删除</el-button
          >
        </el-form-item>
      </el-tab-pane>
      <el-tab-pane label="定时弹幕" name="sixth">
        <el-form-item label="功能开关">
          <el-switch v-model="data.form.CronDanmu" />
        </el-form-item>
        <el-form-item>
          <el-tag>Tips1: 定时弹幕corn表达式 格式为 分 时 日 月 星期 具体参考 https://tool.lu/crontab/ 中的linux格式（此格式支持windows使用）
          </el-tag>
        </el-form-item>
        <el-form-item>
          <el-tag>Tips2: 随即开关打开表示随机发送列表中的一条弹幕, 关闭则是顺序发送</el-tag>
        </el-form-item>
        <el-form-item>
        <el-button
                   type="primary"
                   @click="addDanmulist"
        >新增弹幕列表</el-button
        >
        </el-form-item>
        <div
          v-for="(items, index) in data.form.CronDanmuList"
        >
          <el-form-item :label="'弹幕列表'+ (index+1)">
            <el-button type="danger" @click="deleteDanmulist(index)"
          >删除</el-button
          >
          </el-form-item>
          <el-form-item
          :label="'发送时间'"
        >
          <el-input v-model="data.form.CronDanmuList[index].Cron" />
        </el-form-item>
          <el-form-item
            :label="'随机开关'"
          >
            <el-switch v-model="data.form.CronDanmuList[index].Random" />
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              @click="addDanmu(index)"
            >新增弹幕</el-button
            >
          </el-form-item>

          <el-form-item :label="'弹幕'+(index2+1)" v-for="(items, index2) in data.form.CronDanmuList[index].Danmu">
            <el-input v-model="data.form.CronDanmuList[index].Danmu[index2]" />
            <el-button class="mt-2"
                       @click="deleteDanmu(index,index2)"
                       type="danger"
            >删除</el-button
            >
          </el-form-item>

        </div>
      </el-tab-pane>
      <el-tab-pane label="欢迎黑名单-模糊匹配" name="seventh">
        <el-form-item>
          <el-tag>Tips: 模糊匹配关键字, 只要名字中包含任意一条就不欢迎</el-tag>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="addWelcomeBlacklistWide">新增</el-button>
        </el-form-item>
        <el-form-item
          v-for="(items, index) in data.form.WelcomeBlacklistWide"
          :label="'模糊关键词 ' + (index+1)"
        >
          <el-input v-model="data.form.WelcomeBlacklistWide[index]" />
          <el-button class="mt-2"
                     @click="deleteWelcomeBlacklistWide(index)"
          >删除</el-button
          >
        </el-form-item>
      </el-tab-pane>
      <el-tab-pane label="欢迎黑名单-精确匹配" name="eighth">
        <el-form-item>
          <el-tag>Tips: 精确匹配关键字, 名字完全匹配才不欢迎</el-tag>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="addWelcomeBlacklist">新增</el-button>
        </el-form-item>
        <el-form-item
          v-for="(items, index) in data.form.WelcomeBlacklist"
          :label="'精确关键词 ' + (index+1)"
        >
          <el-input v-model="data.form.WelcomeBlacklist[index]" />
          <el-button class="mt-2"
                     @click="deleteWelcomeBlacklist(index)"
          >删除</el-button
          >
        </el-form-item>
      </el-tab-pane>
      <center>
        <el-button type="success" @click="saveConfig">保存</el-button>
      </center>


    </el-tabs>

  </el-form>
  <el-dialog
    v-model="data.dialogVisible"
    :title="data.dialogtitle"
    width="30%"
  >
   <el-form-item label="用户uid">
     <el-input v-model.number="data.dialogdata.uid"></el-input>
   </el-form-item>
    <el-form-item label="欢迎语">
      <el-input v-model="data.dialogdata.msg"></el-input>
    </el-form-item>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="data.dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveWelcomeString">
          保存
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<style>
.demo-tabs > .el-tabs__content {
  padding: 32px;
  color: #6b778c;
  font-size: 32px;
  font-weight: 600;
}
</style>
