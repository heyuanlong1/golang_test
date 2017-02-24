http://tonybai.com/2014/09/29/a-channel-compendium-for-golang/
http://www.tuicool.com/articles/RVRbyei

for {
        select {
        case x := <- somechan:
            // … 使用x进行一些操作

        case y, ok := <- someOtherchan:
            // … 使用y进行一些操作，
            // 检查ok值判断someOtherchan是否已经关闭

        case outputChan <- z:
            // … z值被成功发送到Channel上时

        default:
            // … 上面case均无法通信时，执行此分支
        }
}