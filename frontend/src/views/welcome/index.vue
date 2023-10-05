<script setup lang="ts">
import { onMounted, reactive } from "vue";

defineOptions({
  name: "Welcome"
});
const formRef = ref<FormInstance>();
const data = reactive({
  isrunning: false,
  savestatus: false,
  savemsg: "",
  tabledata: [],
  dialogdata: {
    uid: "",
    msg: ""
  },
  dialogtitle: "",
  dialogVisible: false,
  dialogAction: false,
  form: {
    RoomId: 3,
    DanmuLen: 20,
    EntryMsg: "花花机器人进入直播间",
    PKNotice: true,
    InteractWord: true,
    EntryEffect: true,
    ThanksGift: true,
    ThanksGiftTimeout: 3,
    WelcomeSwitch: true,
    WelcomeString: {
      "123456": "欢迎宇宙无敌最帅的xxx进入直播间"
    },
    WelcomeDanmu: ["欢迎 {user} ~", "欢迎 {user} 木嘛~", "欢迎 {user} 好诶~"],
    RobotName: "花花",
    TalkRobotCmd: "花花",
    RobotMode: "QingYunKe",
    FuzzyMatchCmd: false,
    ChatGPT: {
      APIToken: "",
      APIUrl: "",
      Prompt: "",
      Limit: false
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
      Random: true,
      Danmu: [
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
    WelcomeBlacklistWide: [],
    WelcomeBlacklist: [],
    InteractWordByTime: false,
    WelcomeDanmuByTime: []
  }
});
const WelcomeDanmuByTimeTemplate = [{
  Enabled: true,
  Key: "earlymorning",
  Random: true,
  Danmu: [
    "欢迎 {user}, 凌晨的问候~",
    "欢迎 {user}, 早安，幸福满满",
    "欢迎 {user}, 祝梦想成真",
    "欢迎 {user}, 望快乐起步",
    "欢迎 {user}, 望奋斗加油",
    "欢迎 {user}, 祝充满力量",
    "欢迎 {user}, 祝精神旺盛",
    "欢迎 {user}, 望谦虚前行",
    "欢迎 {user}, 祝青春无限",
    "欢迎 {user}, 祝前程似锦"
  ]
},
  {
    Enabled: true,
    Key: "morning",
    Random: true,
    Danmu: [
      "欢迎 {user}, 早安，美好开始",
      "欢迎 {user}, 早安，成功在望",
      "欢迎 {user}, 早安，加油！",
      "欢迎 {user}, 早安，勇往直前",
      "欢迎 {user}, 早安，工作愉快",
      "欢迎 {user}, 早安，眉开眼笑",
      "欢迎 {user}, 早安，携手前行",
      "欢迎 {user}, 早安，让生活更好",
      "欢迎 {user}, 早安，热情满怀"
    ]
  },
  {
    Enabled: true,
    Key: "latemorning",
    Random: true,
    Danmu: [
      "欢迎 {user}, 上午好，奋斗有力",
      "欢迎 {user}, 上午好，朝气蓬勃",
      "欢迎 {user}, 上午好，加油向前",
      "欢迎 {user}, 上午好，祝福无限",
      "欢迎 {user}, 上午好，暖心问候",
      "欢迎 {user}, 上午好，快乐工作",
      "欢迎 {user}, 上午好，信心十足",
      "欢迎 {user}, 上午好，开心一整天",
      "欢迎 {user}, 上午好，精神满满",
      "欢迎 {user}, 上午好，工作愉快"
    ]
  },
  {
    Enabled: true,
    Key: "noon",
    Random: true,
    Danmu: [
      "欢迎 {user}, 中午好，午餐愉快",
      "欢迎 {user}, 中午好，美食充盈",
      "欢迎 {user}, 中午好，快乐午时",
      "欢迎 {user}, 中午好，记得午间小憩",
      "欢迎 {user}, 中午好，吃好喝好",
      "欢迎 {user}, 中午好，百般顺利"
    ]
  },
  {
    Enabled: true,
    Key: "afternoon",
    Random: true,
    Danmu: [
      "欢迎 {user}, 下午好，动力十足",
      "欢迎 {user}, 下午好，活力满满",
      "欢迎 {user}, 下午好，祝你成功",
      "欢迎 {user}, 下午好，笑容满面",
      "欢迎 {user}, 下午好，开心每时",
      "欢迎 {user}, 下午好，工作顺利",
      "欢迎 {user}, 下午好，美好心情",
      "欢迎 {user}, 下午好，活力四射",
      "欢迎 {user}, 下午好，幸福无限"
    ]
  },
  {
    Enabled: true,
    Key: "night",
    Random: true,
    Danmu: [
      "欢迎 {user}, 晚上好，祝福相伴",
      "欢迎 {user}, 晚上好，幸福安康",
      "欢迎 {user}, 晚上好，心情愉悦",
      "欢迎 {user}, 晚上好，生活美满",
      "欢迎 {user}, 晚上好，夜色温馨",
      "欢迎 {user}, 晚上好，美梦成真",
      "欢迎 {user}, 晚上好，美好的晚上",
      "欢迎 {user}, 祝愿你今晚开心"
    ]
  },
  {
    Enabled: true,
    Key: "midnight",
    Random: true,
    Danmu: [
      "欢迎 {user}, 午夜好，月色清朗",
      "欢迎 {user}, 午夜好，真诚问候",
      "欢迎 {user}, 午夜好，祝福相随",
      "欢迎 {user}, 午夜深沉，你不是独自一人",
      "欢迎 {user}, 午夜好，静夜祝福",
      "欢迎 {user}, 祝福你坚定向前",
      "欢迎 {user}, 午夜沉寂，吉祥相随。",
      "欢迎 {user}, 午夜好，还没休息?",
      "欢迎 {user}, 拥抱午夜的美好"
    ]
  }
];
//凌晨 - Early morning   2:00--5:00
// 早晨 - Morning   5:00--9:00
// 上午 - Late morning / Mid-morning  9:00--11:00
// 中午 - Noon  11:00--14:00
// 下午 - Afternoon 14:00 -- 20:00
// 晚上 - Evening / Night 20:00--00:00
// 午夜 - Midnight 00:00 -- 2:00
const WelcomeDanmuByTimeDescribe = reactive({
  earlymorning: "凌晨2:00-5:00",
  morning: "早晨5:00-9:00",
  latemorning: "上午9:00-11:00",
  noon: "中午11:00-14:00",
  afternoon: "下午14:00-20:00",
  night: "晚上20:00-00:00",
  midnight: "午夜00:00-2:00"
});
import { ref } from "vue";
import type { FormInstance, FormRules, TabsPaneContext } from "element-plus";
import { ReadConfig, WriteConfig } from "../../../wailsjs/go/main/App";
import { ElNotification } from "element-plus";
import { Monitor, Start, Stop } from "../../../wailsjs/go/main/Program";
// import { func } from "vue-types";

const activeName = ref("first");

const handleClick = (tab: TabsPaneContext, event: Event) => {
  console.log(tab, event);
};
const addWelcomeDanmu = () => {
  data.form.WelcomeDanmu.unshift("{user}");
};
const addWelcomeDanmuByTime = (index) => {
  data.form.WelcomeDanmuByTime[index].Danmu.unshift("{user}");
};
const addWelcomeBlacklist = () => {
  data.form.WelcomeBlacklist.push("");
};
const addWelcomeBlacklistWide = () => {
  data.form.WelcomeBlacklistWide.push("");
};
const deleteWelcomeDanmu = (item: number) => {
  // const index = data.form.WelcomeDanmu.indexOf(item)
  if (item !== -1) {
    data.form.WelcomeDanmu.splice(item, 1);
  }
};
const deleteWelcomeDanmuByTime = (index1,index2) => {
  // const index = data.form.WelcomeDanmu.indexOf(item)
  if (index1 !== -1 && index2 !== -1) {
    data.form.WelcomeDanmuByTime[index1].Danmu.splice(index2, 1);
  }
};
const initWelcomeDanmuByTime = () => {
  if (data.form.WelcomeDanmuByTime == null){
    data.form.WelcomeDanmuByTime = WelcomeDanmuByTimeTemplate;
    return;
  }
  if (data.form.WelcomeDanmuByTime.length == 0) {
    data.form.WelcomeDanmuByTime = WelcomeDanmuByTimeTemplate;
    return;
  }
  if (data.form.WelcomeDanmuByTime.length < 7) {
    for (let i = 0; i < WelcomeDanmuByTimeTemplate.length; i++) {
      let x = true;
      for (let a = 0; i < data.form.WelcomeDanmuByTime.length; i++) {
        if (data.form.WelcomeDanmuByTime[a].Key === WelcomeDanmuByTimeTemplate[i].Key) {
          x = false;
          break;
        }
      }
      if (x) {
        data.form.WelcomeDanmuByTime.push(WelcomeDanmuByTimeTemplate[i]);
      }
    }
  }
  return;
};
const deleteWelcomeBlacklist = (item: number) => {
  // const index = data.form.WelcomeDanmu.indexOf(item)
  if (item !== -1) {
    data.form.WelcomeBlacklist.splice(item, 1);
  }
};
const deleteWelcomeBlacklistWide = (item: number) => {
  // const index = data.form.WelcomeDanmu.indexOf(item)
  if (item !== -1) {
    data.form.WelcomeBlacklistWide.splice(item, 1);
  }
};
const addFocusDanmu = () => {
  data.form.FocusDanmu.push("");
};
const deleteFocusDanmu = (item: number) => {
  // const index = data.form.WelcomeDanmu.indexOf(item)
  if (item !== -1) {
    data.form.FocusDanmu.splice(item, 1);
  }
};
const getWelcomeDanmuByTimeDescribe = (key) => {
  return WelcomeDanmuByTimeDescribe[key];
};

function formatWelcomeString() {
  const arr = Object.entries(data.form.WelcomeString);
  data.tabledata = arr.map(([key, value]) => ({ key, value }));
}

function addWelcomeString() {
  data.dialogdata = {
    uid: "",
    msg: ""
  };
  data.dialogtitle = "新增";
  data.dialogAction = false;
  data.dialogVisible = true;

}

function saveWelcomeString() {
  data.form.WelcomeString[data.dialogdata.uid] = data.dialogdata.msg;
  formatWelcomeString();
  data.dialogVisible = false;
}

function deleteWelcomeString(item: string) {
  // console.log(item)
  if (data.form.WelcomeString.hasOwnProperty(item)) {
    delete data.form.WelcomeString[item];
  }
  formatWelcomeString();
}

function editWelcomeString(uid, msg) {
  data.dialogdata.uid = uid;
  data.dialogdata.msg = msg;
  data.dialogAction = true;
  data.dialogVisible = true;
}

function addDanmulist() {
  data.form.CronDanmuList.push({
    Cron: "*/2 * * * *",
    Random: true,
    Danmu: [
      ""
    ]
  });
}

function deleteDanmulist(row: number) {
  if (row !== -1) {
    data.form.CronDanmuList.splice(row, 1);
  }
}

function addDanmu(item: number) {
  data.form.CronDanmuList[item].Danmu.push("");
}

function deleteDanmu(item: number, row: number) {
  if (item !== -1 && row !== -1) {
    data.form.CronDanmuList[item].Danmu.splice(row, 1);
  }
}

async function saveConfig() {
  data.savestatus = false;
  data.savemsg = "";
  console.log(data.form.ChatGPT.APIUrl)
  if (data.form.ChatGPT.APIUrl.length == 0) {
    data.form.ChatGPT.APIUrl = "https://api.openai.com/v1";
  }
  await WriteConfig(JSON.stringify(data.form)).then(res => {
    data.savestatus = res.Code;
    data.savemsg = res.Msg;
  });
  if (data.savestatus) {
    ElNotification({
      title: "操作成功",
      message: "操作成功，正在重启机器人",
      type: "success"
    });
    restart();
  } else {
    ElNotification({
      title: "保存失败",
      message: "请修改配置直到保存正确为止，错误信息：" + data.savemsg,
      type: "error"
    });
  }
}

async function pgstart() {
  await Monitor().then(res => {
    data.isrunning = res;
  });
  if (data.isrunning == false) {
    await Start().then(res => {
      data.isrunning = res;
    });
  }
}

async function pgstop() {
  await Monitor().then(res => {
    data.isrunning = res;
  });
  if (data.isrunning == true) {
    await Stop().then(res => {
      data.isrunning = !res;
    });
  }
}

async function restart() {
  await pgstop();
  console.log(data.isrunning);
  if (data.isrunning == false) {
    pgstart();
  }
}

const rules = reactive<FormRules>({
  RoomId: [
    { required: true, message: "此项必填" },
    { type: "number", message: "此处必须为数字" }
  ],
  DanmuLen: [
    { required: true, message: "此项必填" },
    { type: "number", message: "此处必须为数字" }
  ]
});
onMounted(() => {
  ReadConfig().then(res => {
    if (!res.Code) {
      console.log(res);
      ElNotification({
        title: "读取配置文件失败",
        message: res.Msg,
        type: "error"
      });
    } else {
      data.form = res.Form;
      formatWelcomeString();
      initWelcomeDanmuByTime();
      if (data.form.ChatGPT.APIUrl.length == 0) {
        data.form.ChatGPT.APIUrl = "https://api.openai.com/v1";
      }
    }
  });
});
</script>

<template>
  <el-form :model="data.form" label-width="120px" ref="formRef" class="el-form" :rules="rules">
    <el-tabs v-model="activeName" class="demo-tabs" type="border-card" @tab-click="handleClick">
      <el-tab-pane label="基础配置" name="first">
        <el-form-item label="直播间号" prop="RoomId">
          <el-input v-model.number="data.form.RoomId" autocomplete="off" />
        </el-form-item>
        <el-form-item label="可发送弹幕长度" prop="DanmuLen">
          <el-input v-model.number="data.form.DanmuLen" />
        </el-form-item>
        <el-form-item label="进入直播间消息" prop="DanmuLen">
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
        <el-form-item label="礼物感谢频率" prop="RoomId">
          <el-col :span="3">
          <el-input v-model.number="data.form.ThanksGiftTimeout" >
            <template #append>
              <div class="input-append">秒</div>
            </template>
          </el-input>
          </el-col>
        </el-form-item>
        <el-form-item label="PK提醒">
          <el-switch v-model="data.form.PKNotice" />
        </el-form-item>

      </el-tab-pane>

      <el-tab-pane label="欢迎弹幕自定义" name="second">

        <el-form-item>
          <el-tag>Tips: 自定义欢迎语 {user}为用户昵称占位符，列表为随机列表</el-tag>
        </el-form-item>
        <el-form-item label="时间段时间欢迎">
          <el-switch v-model="data.form.InteractWordByTime" />
        </el-form-item>
        <el-form-item v-if="data.form.InteractWordByTime == false">
          <el-button type="primary" @click="addWelcomeDanmu">新增</el-button>
        </el-form-item>
        <el-form-item
          v-if="data.form.InteractWordByTime == false"
          v-for="(items, index) in data.form.WelcomeDanmu"
          :label="'欢迎弹幕 ' + (index+1)"
        >
          <el-input v-model="data.form.WelcomeDanmu[index]" />
          <el-button class="mt-2" @click="deleteWelcomeDanmu(index)"
          >删除
          </el-button>
        </el-form-item>
        <div v-if="data.form.InteractWordByTime == true" v-for="(items, index) in data.form.WelcomeDanmuByTime">
          <el-text class="el-text--large" type="primary">{{getWelcomeDanmuByTimeDescribe(data.form.WelcomeDanmuByTime[index].Key)}}</el-text>
          <el-form-item>
            <el-button type="primary" @click="addWelcomeDanmuByTime(index)">新增</el-button>
          </el-form-item>
          <el-form-item label="启用">
            <el-switch v-model="data.form.WelcomeDanmuByTime[index].Enabled" />
          </el-form-item>
          <el-form-item label="随机弹幕">
            <el-switch v-model="data.form.WelcomeDanmuByTime[index].Random" />
          </el-form-item>
          <el-form-item
            v-for="(items, index2) in data.form.WelcomeDanmuByTime[index].Danmu"
            :label="'欢迎弹幕 ' + (index2+1)"
          >
            <el-input v-model="data.form.WelcomeDanmuByTime[index].Danmu[index2]" />
            <el-button class="mt-2" @click="deleteWelcomeDanmuByTime(index,index2)"
            >删除
            </el-button>
          </el-form-item>

        </div>

      </el-tab-pane>
      <el-tab-pane label="指定人欢迎" name="third">
        <!--        <template>-->
        <el-form-item label="功能开关">
          <el-switch v-model="data.form.WelcomeSwitch" />
        </el-form-item>
        <el-button @click="addWelcomeString" type="primary">添加</el-button>
        <el-table :data="data.tabledata" border>
          <el-table-column prop="key" label="用户uid" />
          <el-table-column prop="value" label="欢迎语" />
          <el-table-column label="操作">
            <template #default="scope">
              <el-button type="primary" @click="editWelcomeString(scope.row.key,scope.row.value)">编辑</el-button>
              <el-button type="danger" @click="deleteWelcomeString(scope.row.key)">删除</el-button>
            </template>
          </el-table-column>
          <!--            <el-table-column prop="address" label="Address" />-->
        </el-table>
        <!--        </template>-->
      </el-tab-pane>
      <el-tab-pane label="ai机器人" name="fourth">
        <el-form-item>
          <el-tag>Tips: 直播间调用机器人的昵称模糊匹配，"花花今天你吃了吗"与"今天花花你吃了吗"都会触发机器人回答</el-tag>
        </el-form-item>
        <el-form-item label="机器人名称">
          <el-input v-model.number="data.form.RobotName" />
        </el-form-item>
        <el-form-item label="触发关键字">
          <el-input v-model.number="data.form.TalkRobotCmd" />
        </el-form-item>
        <el-form-item label="关键词模糊匹配">
          <el-switch v-model="data.form.FuzzyMatchCmd" />
        </el-form-item>
        <el-form-item label="服务提供商" prop="region">
          <el-select v-model="data.form.RobotMode" placeholder="请选择">
            <el-option label="QingYunKe" value="QingYunKe" />
            <el-option label="ChatGPT" value="ChatGPT" />
          </el-select>
        </el-form-item>
        <el-form-item label="API KEY" v-if="data.form.RobotMode == 'ChatGPT'">
          <el-input v-model="data.form.ChatGPT.APIToken" />
        </el-form-item>
        <el-form-item label="API URL" v-if="data.form.RobotMode == 'ChatGPT'">
          <el-input v-model="data.form.ChatGPT.APIUrl" />
        </el-form-item>
        <el-form-item label="prompt" v-if="data.form.RobotMode == 'ChatGPT'">
          <el-input v-model="data.form.ChatGPT.Prompt" />
        </el-form-item>
        <el-form-item label="弹幕长度限制" v-if="data.form.RobotMode == 'ChatGPT'">
          <el-switch v-model="data.form.ChatGPT.Limit" />
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
          >删除
          </el-button
          >
        </el-form-item>
      </el-tab-pane>
      <el-tab-pane label="定时弹幕" name="sixth">
        <el-form-item label="功能开关">
          <el-switch v-model="data.form.CronDanmu" />
        </el-form-item>
        <el-form-item>
          <el-tag>Tips1: 定时弹幕corn表达式 格式为 秒 分 时 日 月 星期 其中秒是可选的 具体参考 https://tool.lu/crontab/
            中的linux格式（此格式支持windows使用）
          </el-tag>

        </el-form-item>
        <el-form-item>
          <el-tag>
            Tips2: 想要用秒的时候, 应该是 30 * * * * * * 共六项 表示每分钟的第30秒执行
          </el-tag>
        </el-form-item>
        <el-form-item>
          <el-tag>
          Tips3: 不想用秒的时候, 应该是 */1 * * * * 共五项 表示每1分钟执行一次
          </el-tag>
        </el-form-item>
        <el-form-item>
          <el-tag>Tips4: 随机开关打开表示随机发送列表中的一条弹幕, 关闭则是顺序发送</el-tag>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            @click="addDanmulist"
          >新增弹幕列表
          </el-button
          >
        </el-form-item>
        <div
          v-for="(items, index) in data.form.CronDanmuList"
        >
          <el-form-item :label="'弹幕列表'+ (index+1)">
            <el-button type="danger" @click="deleteDanmulist(index)"
            >删除
            </el-button
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
            >新增弹幕
            </el-button
            >
          </el-form-item>

          <el-form-item :label="'弹幕'+(index2+1)" v-for="(items, index2) in data.form.CronDanmuList[index].Danmu">
            <el-input v-model="data.form.CronDanmuList[index].Danmu[index2]" />
            <el-button class="mt-2"
                       @click="deleteDanmu(index,index2)"
                       type="danger"
            >删除
            </el-button
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
          >删除
          </el-button
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
          >删除
          </el-button
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
      <el-input v-model.number="data.dialogdata.uid" :disabled=data.dialogAction></el-input>
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
.el-tabs__header{
  position: sticky;
    margin: 0 0 15px;
    top: 0;
    z-index: 2;
}
</style>
