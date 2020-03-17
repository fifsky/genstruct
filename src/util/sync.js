import { Err } from './error'
import { notification } from 'antd'

export const sync = (fn,isCatch = true) => fn().catch(function (e) {
  if(isCatch) {
    notification.error({
      message: `错误`,
      description: Err.instance(e).getMsg(),
    })
  }
})
