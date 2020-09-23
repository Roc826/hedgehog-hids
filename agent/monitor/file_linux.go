package monitor

import (
	"github.com/fsnotify/fsnotify"
	"hedgehog-hids-agent/common"
	"hedgehog-hids-agent/config"
	"hedgehog-hids-agent/log"
	"os"
	"path"
	"path/filepath"
)

func FileMonitor(resultChan chan map[string]string) {
	conf := config.LoadConfig(config.ConfigPath);
	watcher,err := fsnotify.NewWatcher()
	if err != nil{
		log.Error(err)
	}
	for _,dir :=range(conf.MonitorDirs){
		err := 	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				path, err := filepath.Abs(path)
				if err != nil {
					log.Error(err)
					return nil
				}

				err = watcher.Add(path)
				if err != nil {
					log.Error(err)
					return nil
				}
				log.Debug("添加监控:", path)
			}
			return nil
		})
		if err != nil {
			log.Error(err)
			return
		}
	}

	go func(){
		for {
			select {
			case event:= <-watcher.Events:{
				//获得文件尾缀
				nameExtension := path.Ext(event.Name)
				//去除文件尾缀的第一个字符.
				if len(nameExtension) != 0 {
					nameExtension = nameExtension[1:len(nameExtension)]
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Debug("Create:", event.Name)
					//这里获取创建文件的信息,如果是在目录中,则加入监控中
					//如果有nameExtension == ".jpg"限制,则无法创建文件夹
					fi, err := os.Stat(event.Name)
					if len(conf.WhiteExt) != 0 {
						if common.InSlice(nameExtension, conf.WhiteExt) {
							if err == nil && fi.IsDir() {
								err := watcher.Add(event.Name)
								log.Debug("添加监控:", event.Name)
								if err != nil{
									log.Error(err)
								}
							}
						}else{
							log.Debug("文件未在白名单中,移除文件",event.Name)
							err := os.Remove(event.Name)
							if err != nil{
								log.Error("移除文件失败:",err)
							}
						}
						continue
					}
					if len(conf.BlackExt) != 0 && common.InSlice(nameExtension,conf.BlackExt){
						log.Debug("文件在黑名单中,移除文件",event.Name)
						err := os.Remove(event.Name)
						if err != nil{
							log.Error("移除文件失败:",err)
						}
					}
					if err == nil && fi.IsDir() {
						watcher.Add(event.Name)
						log.Debug("添加监控:", event.Name)
					}
				}
				//write
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Debug("Write:", event.Name)
				}
				//remove
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Debug("Delete:", event.Name)
					watcher.Remove(event.Name)

				}
				//rename 重命名
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Debug("Rename:", event.Name)
					//如果重命名文件是目录,则移除监控
					//因为文件如果被移除了也无法使用os.Stat判断是不是目录,所以无脑移除
					watcher.Remove(event.Name)
				}
				//chmod
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Debug("修改权限:", event.Name)
				}
			}
			case err := <-watcher.Errors:{
				log.Error(err)
			}

			}
		}
	}()

}
