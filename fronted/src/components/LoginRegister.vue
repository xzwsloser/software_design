<template>
  <div class="login-register-container">
    <div class="form-card">
      <!-- 标题 -->
      <div class="title">
        <h2>可视化景点推荐系统</h2>
      </div>

      <!-- Tab切换 -->
      <el-tabs v-model="activeTab" class="tabs" @tab-change="handleTabChange">
        <el-tab-pane label="登录" name="login">
          <el-form
            :model="loginForm"
            label-width="0"
            class="form-content"
          >
            <el-form-item>
              <el-input
                v-model="loginForm.username"
                placeholder="用户名"
                size="large"
              >
                <template #prefix>
                  <el-icon><User /></el-icon>
                </template> 
              </el-input>

            </el-form-item>
            <el-form-item>
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="密码"
                size="large"
                show-password
              >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                size="large"
                class="submit-btn"
                @click="handleLogin"
              >
                登录
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="注册" name="register">
          <el-form
            :model="registerForm"
            label-width="0"
            class="form-content"
          >
            <el-form-item>
              <el-input
                v-model="registerForm.username"
                placeholder="用户名"
                size="large"
              >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
            </el-form-item>

            <el-form-item>
              <el-select
                v-model="registerForm.gender"
                placeholder="请选择性别"
                size="large"
                style="width: 100%"
              >
                <el-option
                  label="男"
                  :value="0"
                />
                <el-option
                  label="女"
                  :value="1"
                />
              </el-select>
            </el-form-item>

            <el-form-item>
              <el-select
                v-model="selectedProvince"
                placeholder="请选择省份"
                size="large"
                style="width: 100%"
                @change="handleProvinceChange"
              >
                <el-option
                  v-for="province in cities"
                  :key="province.name"
                  :label="province.name"
                  :value="province.name"
                />
              </el-select>
            </el-form-item>

            <el-form-item v-if="selectedProvince && currentCities.length !== 0">
              <el-select
                v-model="registerForm.city"
                placeholder="请选择城市"
                size="large"
                style="width: 100%"
              >
                <el-option
                  v-for="city in currentCities"
                  :key="city"
                  :label="city"
                  :value="city"
                />
              </el-select>
            </el-form-item>

            <el-form-item>
              <el-input
                v-model="registerForm.password"
                type="password"
                placeholder="密码"
                size="large"
                show-password
              >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
            </el-form-item>

            <el-form-item>
              <el-input
                v-model="registerForm.confirmPassword"
                type="password"
                placeholder="确认密码"
                size="large"
                show-password
              >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                size="large"
                class="submit-btn"
                @click="handleRegister"
              >
                注册
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../stores/user'
import { User, Lock } from '@element-plus/icons-vue'

// Tab切换
const activeTab = ref('login')
const isLogin = computed(() => activeTab.value === 'login')

// 使用路由和用户store
const router = useRouter()
const userStore = useUserStore()

// 登录表单
const loginForm = reactive({
  username: '',
  password: '',
})

// 选中的省份
const selectedProvince = ref('')


// 中国城市数据
const cities = [
  {
    name: '北京市',
    // cities: ['东城区', '西城区', '朝阳区', '丰台区', '石景山区', '海淀区', '门头沟区', '房山区', '通州区', '顺义区', '昌平区', '大兴区', '怀柔区', '平谷区', '密云区', '延庆区']
    cities: []
  },
  {
    name: '天津市',
    // cities: ['和平区', '河东区', '河西区', '南开区', '河北区', '红桥区', '东丽区', '西青区', '津南区', '北辰区', '武清区', '宝坻区', '滨海新区', '宁河区', '静海区', '蓟州区']
    cities: []
  },
  {
    name: '河北省',
    cities: ['石家庄市', '唐山市', '秦皇岛市', '邯郸市', '邢台市', '保定市', '张家口市', '承德市', '沧州市', '廊坊市', '衡水市']
  },
  {
    name: '山西省',
    cities: ['太原市', '大同市', '阳泉市', '长治市', '晋城市', '朔州市', '晋中市', '运城市', '忻州市', '临汾市', '吕梁市']
  },
  {
    name: '内蒙古自治区',
    cities: ['呼和浩特市', '包头市', '乌海市', '赤峰市', '通辽市', '鄂尔多斯市', '呼伦贝尔市', '巴彦淖尔市', '乌兰察布市']
  },
  {
    name: '辽宁省',
    cities: ['沈阳市', '大连市', '鞍山市', '抚顺市', '本溪市', '丹东市', '锦州市', '营口市', '阜新市', '辽阳市', '盘锦市', '铁岭市', '朝阳市', '葫芦岛市']
  },
  {
    name: '吉林省',
    cities: ['长春市', '吉林市', '四平市', '辽源市', '通化市', '白山市', '松原市', '白城市', '延边朝鲜族自治州']
  },
  {
    name: '黑龙江省',
    cities: ['哈尔滨市', '齐齐哈尔市', '鸡西市', '鹤岗市', '双鸭山市', '大庆市', '伊春市', '佳木斯市', '七台河市', '牡丹江市', '黑河市', '绥化市', '大兴安岭地区']
  },
  {
    name: '上海市',
    // cities: ['黄浦区', '徐汇区', '长宁区', '静安区', '普陀区', '虹口区', '杨浦区', '闵行区', '宝山区', '嘉定区', '浦东新区', '金山区', '松江区', '青浦区', '奉贤区', '崇明区']
    cities: []
  },
  {
    name: '江苏省',
    cities: ['南京市', '无锡市', '徐州市', '常州市', '苏州市', '南通市', '连云港市', '淮安市', '盐城市', '扬州市', '镇江市', '泰州市', '宿迁市']
  },
  {
    name: '浙江省',
    cities: ['杭州市', '宁波市', '温州市', '嘉兴市', '湖州市', '绍兴市', '金华市', '衢州市', '舟山市', '台州市', '丽水市']
  },
  {
    name: '安徽省',
    cities: ['合肥市', '芜湖市', '蚌埠市', '淮南市', '马鞍山市', '淮北市', '铜陵市', '安庆市', '黄山市', '滁州市', '阜阳市', '宿州市', '六安市', '亳州市', '池州市', '宣城市']
  },
  {
    name: '福建省',
    cities: ['福州市', '厦门市', '莆田市', '三明市', '泉州市', '漳州市', '南平市', '龙岩市', '宁德市']
  },
  {
    name: '江西省',
    cities: ['南昌市', '景德镇市', '萍乡市', '九江市', '新余市', '鹰潭市', '赣州市', '吉安市', '宜春市', '抚州市', '上饶市']
  },
  {
    name: '山东省',
    cities: ['济南市', '青岛市', '淄博市', '枣庄市', '东营市', '烟台市', '潍坊市', '济宁市', '泰安市', '威海市', '日照市', '临沂市', '德州市', '聊城市', '滨州市', '菏泽市']
  },
  {
    name: '河南省',
    cities: ['郑州市', '开封市', '洛阳市', '平顶山市', '安阳市', '鹤壁市', '新乡市', '焦作市', '濮阳市', '许昌市', '漯河市', '三门峡市', '南阳市', '商丘市', '信阳市', '周口市', '驻马店市']
  },
  {
    name: '湖北省',
    cities: ['武汉市', '黄石市', '十堰市', '宜昌市', '襄阳市', '鄂州市', '荆门市', '孝感市', '荆州市', '黄冈市', '咸宁市', '随州市', '恩施土家族苗族自治州']
  },
  {
    name: '湖南省',
    cities: ['长沙市', '株洲市', '湘潭市', '衡阳市', '邵阳市', '岳阳市', '常德市', '张家界市', '益阳市', '郴州市', '永州市', '怀化市', '娄底市', '湘西土家族苗族自治州']
  },
  {
    name: '广东省',
    cities: ['广州市', '韶关市', '深圳市', '珠海市', '汕头市', '佛山市', '江门市', '湛江市', '茂名市', '肇庆市', '惠州市', '梅州市', '汕尾市', '河源市', '阳江市', '清远市', '东莞市', '中山市', '潮州市', '揭阳市', '云浮市']
  },
  {
    name: '广西壮族自治区',
    cities: ['南宁市', '柳州市', '桂林市', '梧州市', '北海市', '防城港市', '钦州市', '贵港市', '玉林市', '百色市', '贺州市', '河池市', '来宾市', '崇左市']
  },
  {
    name: '海南省',
    cities: ['海口市', '三亚市', '三沙市', '儋州市', '五指山市', '琼海市', '文昌市', '万宁市', '东方市', '定安县', '屯昌县', '澄迈县', '临高县', '白沙黎族自治县', '昌江黎族自治县', '乐东黎族自治县', '陵水黎族自治县', '保亭黎族苗族自治县', '琼中黎族苗族自治县']
  },
  {
    name: '重庆市',
    // cities: ['渝中区', '万州区', '涪陵区', '大渡口区', '江北区', '沙坪坝区', '九龙坡区', '南岸区', '北碚区', '綦江区', '大足区', '渝北区', '巴南区', '黔江区', '长寿区', '江津区', '合川区', '永川区', '南川区', '璧山区', '铜梁区', '潼南区', '荣昌区', '开州区', '梁平区', '武隆区']
    cities: []

  },
  {
    name: '四川省',
    cities: ['成都市', '自贡市', '攀枝花市', '泸州市', '德阳市', '绵阳市', '广元市', '遂宁市', '内江市', '乐山市', '南充市', '眉山市', '宜宾市', '广安市', '达州市', '雅安市', '巴中市', '资阳市', '阿坝藏族羌族自治州', '甘孜藏族自治州', '凉山彝族自治州']
  },
  {
    name: '贵州省',
    cities: ['贵阳市', '六盘水市', '遵义市', '安顺市', '毕节市', '铜仁市', '黔西南布依族苗族自治州', '黔东南苗族侗族自治州', '黔南布依族苗族自治州']
  },
  {
    name: '云南省',
    cities: ['昆明市', '曲靖市', '玉溪市', '保山市', '昭通市', '丽江市', '普洱市', '临沧市', '楚雄彝族自治州', '红河哈尼族彝族自治州', '文山壮族苗族自治州', '西双版纳傣族自治州', '大理白族自治州', '德宏傣族景颇族自治州', '怒江傈僳族自治州', '迪庆藏族自治州']
  },
  {
    name: '西藏自治区',
    cities: ['拉萨市', '日喀则市', '昌都市', '林芝市', '山南市', '那曲市', '阿里地区']
  },
  {
    name: '陕西省',
    cities: ['西安市', '铜川市', '宝鸡市', '咸阳市', '渭南市', '延安市', '汉中市', '榆林市', '安康市', '商洛市']
  },
  {
    name: '甘肃省',
    cities: ['兰州市', '嘉峪关市', '金昌市', '白银市', '天水市', '武威市', '张掖市', '平凉市', '酒泉市', '庆阳市', '定西市', '陇南市', '临夏回族自治州', '甘南藏族自治州']
  },
  {
    name: '青海省',
    cities: ['西宁市', '海东市', '海北藏族自治州', '黄南藏族自治州', '海南藏族自治州', '果洛藏族自治州', '玉树藏族自治州', '海西蒙古族藏族自治州']
  },
  {
    name: '宁夏回族自治区',
    cities: ['银川市', '石嘴山市', '吴忠市', '固原市', '中卫市']
  },
  {
    name: '新疆维吾尔自治区',
    cities: ['乌鲁木齐市', '克拉玛依市', '吐鲁番市', '哈密市', '昌吉回族自治州', '博尔塔拉蒙古自治州', '巴音郭楞蒙古自治州', '阿克苏地区', '克孜勒苏柯尔克孜自治州', '喀什地区', '和田地区', '伊犁哈萨克自治州', '塔城地区', '阿勒泰地区']
  }
]

const specialCities = new Set(['北京市', '天津市', '重庆市', '上海市' ])

// 注册表单
const registerForm = reactive({
  username: '',
  gender: '',
  city: '',
  password: '',
  confirmPassword: ''
})

// 选中的地区
// const areaName = computed(() => {
//   if (!selectedProvince.value && !registerForm.city) return ''
//   if (!registerForm.city) return selectedProvince
//   return `${selectedProvince}${registerForm.city}`
// })

// 当前选中的省份对应的城市列表
const currentCities = computed(() => {
  if (!selectedProvince.value) return []
  const province = cities.find(p => p.name === selectedProvince.value)
  return province ? province.cities : []
})

// Tab切换处理
const handleTabChange = (tabName) => {
  activeTab.value = tabName
}

// 省份变化处理
const handleProvinceChange = () => {
  registerForm.city = ''
}

// 登录处理
const handleLogin = async () => {
  // 表单验证
  if (!loginForm.username || !loginForm.password) {
    ElMessage.error('请填写用户名和密码')
    return
  }

  try {
    const result = await userStore.login(loginForm)
    if (result.success) {
      ElMessage.success('登录成功！')
      // 跳转到景点列表页面
      router.push('/sites')
    } else {
      ElMessage.error(result.error || '登录失败')
    }
  } catch (error) {
    ElMessage.error('登录失败，请稍后重试')
  }
}

// 注册处理
const handleRegister = async () => {
  // 表单验证
  if (!registerForm.username || !registerForm.password) {
    ElMessage.error('请填写用户名和密码')
    return
  }

  if (registerForm.password !== registerForm.confirmPassword) {
    ElMessage.error('两次输入的密码不一致')
    return
  }

  if (registerForm.gender === '') {
    ElMessage.error('请选择性别')
    return
  }

  if (!selectedProvince.value || !registerForm.city) {
    if (!specialCities.has(selectedProvince.value)) {
      ElMessage.error('请选择省份和城市')
      return
    }
  }


  // 添加省份信息到注册表单
  registerForm.province = selectedProvince.value

  try {
    const result = await userStore.register(registerForm)
    if (result.success) {
      ElMessage.success('注册成功！')
      // 跳转到景点列表页面
      router.push('/sites')
    } else {
      ElMessage.error(result.error || '注册失败')
    }
  } catch (error) {
    ElMessage.error('注册失败，请稍后重试')
  }
}
</script>

<style scoped>
.login-register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.form-card {
  width: 100%;
  max-width: 400px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  overflow: hidden;
}

.title {
  text-align: center;
  padding: 30px 0 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.title h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 500;
}

.tabs {
  padding: 0 30px 30px;
}

.tabs :deep(.el-tabs__header) {
  margin-bottom: 30px;
}

.tabs :deep(.el-tabs__nav-wrap::after) {
  display: none;
}

.tabs :deep(.el-tabs__item) {
  padding: 0 30px;
  font-size: 16px;
  font-weight: 500;
}

.tabs :deep(.el-tabs__item.is-active) {
  color: #667eea;
}

.tabs :deep(.el-tabs__active-bar) {
  background-color: #667eea;
}

.form-content {
  width: 100%;
}

.form-content :deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 0 0 1px #dcdfe6 inset;
  transition: all 0.3s;
}

.form-content :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #c0c4cc inset;
}

.form-content :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2) inset;
}


.submit-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  margin-top: 10px;
}

.submit-btn:hover {
  background: linear-gradient(135deg, #5a72d1 0%, #6a4190 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.submit-btn:active {
  transform: translateY(0);
}
</style>