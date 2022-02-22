package log

//用于输出日志
import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

// Settings 存储logger的配置信息，为func Setup提供配置
type Settings struct {
	//文件所在路径
	Path string
	//文件名
	Name string
	//文件后缀
	Ext string
	//时间格式
	TimeFormat string
	//文件夹内文件最大数量
	MaxBackups uint32
	//文件的大小(k)
	MaxSize uint32
	//用于计数到第几个文件
	Cnt int
}

var (
	logger        *log.Logger
	defaultPrefix = ""
	mu            sync.Mutex
	logPrefix     = ""
	levelFlags    = []string{"ERROR", "FATAL"}
	flags         = log.LstdFlags | log.Lmicroseconds | log.Lshortfile
	filePCB       *os.File
	config        *Settings
	fileStat      os.FileInfo
)

type logLevel int

// log levels
const (
	ERROR logLevel = iota
	FATAL
)

func init() {
	logger = log.New(os.Stdout, defaultPrefix, flags)
}

//检查目录是否存在，true为不存在
func checkNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

//检查目录是否有权限，true为没权限
func checkNotPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

func isNotExistMkDir(src string) error {
	if notExist := checkNotExist(src); notExist {
		if err := mkDir(src); err != nil {
			return err
		}
	}
	return nil
}

func mkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm) //Unix permission bits, 0o777
	if err != nil {
		return err
	}
	return nil
}

func mustOpen(fileName, dir string) (*os.File, error) {
	perm := checkNotPermission(dir)
	if perm {
		return nil, fmt.Errorf("permission denied dir: %s", dir)
	}

	err := isNotExistMkDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error during make dir %s, err: %s", dir, err)
	}

	f, err := os.OpenFile(dir+string(os.PathSeparator)+fileName, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to open file, err: %s", err)
	}

	return f, nil
}

// Setup 配置logger
func Setup(settings *Settings) {
	config = settings
	if settings.Cnt == 11 {
		settings.Cnt = 1
	}
	if filePCB != nil {
		filePCB.Close()
		filePCB = nil
	}
	dir := settings.Path + time.Now().Format(settings.TimeFormat)
	fileName := fmt.Sprintf("%s_%s.%s",
		settings.Name,
		strconv.Itoa(settings.Cnt),
		settings.Ext)
	logFile, err := mustOpen(fileName, dir)
	if err != nil {
		log.Fatalf("logging.Setup err: %s", err)
	}
	filePCB = logFile
	mw := io.MultiWriter(os.Stdout, logFile)
	logger = log.New(mw, defaultPrefix, flags)
	settings.Cnt++
}

//写入文件时检查文件大小和日期(mustOpen中其实已经会检查日期了)
func checkFileSize() {
	if filePCB == nil {
		return
	}
	fileStat, _ = filePCB.Stat()
	if fileStat.Size() > int64(config.MaxSize*1024) {
		Setup(config)
	}
}

//设置logger的前缀
func setPrefix(level logLevel) {

	logPrefix = fmt.Sprintf("[%s] ", levelFlags[level])

	logger.SetPrefix(logPrefix)
}

func Error(v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	setPrefix(ERROR)
	checkFileSize()
	logger.Println(v...)
}

// Fatal 打印错误日志，然后停止程序
func Fatal(v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	setPrefix(FATAL)
	checkFileSize()
	logger.Fatalln(v...)
}
