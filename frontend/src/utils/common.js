// 格式化朝代显示
export const formatDynasty = (dynasty) => {
  if (!dynasty) return ''
  
  const dynastyMap = {
    'tang': '唐',
    'song': '宋',
    'yuan': '元',
    'ming': '明',
    'qing': '清',
    'preqin': '先秦',
    'wudai': '五代',
    'nan': '南唐', // 针对五代南唐的特殊处理
    'jin': '金',
    'liao': '辽'
  }
  
  // 转换小写进行匹配，兼容大小写不一致的情况
  const key = dynasty.toLowerCase()
  return dynastyMap[key] || dynasty
}
