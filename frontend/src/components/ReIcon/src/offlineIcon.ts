import { addIcon } from "@iconify/vue/dist/offline";

/**
 * 这里存放本地图标，在 src/layout/index.vue 文件中加载，避免在首启动加载
 */

// 本地菜单图标，后端在路由的icon中返回对应的图标字符串并且前端在此处使用addIcon添加即可渲染菜单图标
import InformationLine from "@iconify-icons/ri/information-line";
import Lollipop from "@iconify-icons/ep/lollipop";
import HomeFilled from "@iconify-icons/ep/home-filled";
// 注意：模块路径必须与包内文件名大小写完全一致（notification.js/sunny.js），
// Linux 构建环境区分大小写，写成 Notification/Sunny 会导致 rollup 解析失败
import Notification from "@iconify-icons/ep/notification";
import Sunny from "@iconify-icons/ep/sunny";

addIcon("Sunny", Sunny);
addIcon("notification", Notification);
addIcon("homeFilled", HomeFilled);
addIcon("lollipop", Lollipop);
addIcon("informationLine", InformationLine);
