//go:generate goversioninfo -icon=manifest/rc.ico -manifest=manifest/main.exe.manifest
package main

import (
	"agent/global"
	"agent/routes/file"
	"agent/routes/health"
	"agent/routes/shell"
	"agent/routes/update"
	"agent/util"
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	m            string
	buildVersion string
)

func init() {
	flag.StringVar(&m, "m", "release", "mode for run")
}

func main() {
	println(`                                                  
                                      r7                                         7r                                                    
                                    :12USr                                     ;SS22:                                                  
                                  :sULsUU;                                     :2UsLUJ:                                                
                                ,s2sLUJ:                                         :sULs2J,                                              
                              .vSss12:                                             :wwssUc,                                            
                             721c11i                                                 ;12LJUc.                                          
                           rUJLJUr                                                     ;UJLsSr.                                        
                         rw1cJ27                                                         7UJc1wr                                       
                       :21LsSv.                                                           .75ss1U;                                     
                     :J2Lsws.                                                               .LwJLUJ;                                   
                   ,JwLvJ7,                                                                   .7J7Lw1:                                 
                 ,LUL777.                                                                       .777cUs:                               
                ,Ps77r77rrccLcLLsLsLsLsLsLsLsLsLsLsLsLsLsLsLsLsLsLsLsLsLsLsLsLsLsLsLsLsLsLLcsccr;r7r77sK:                              
                .577r777vLcLcscLLLLscLLsLscsLLcscsLLcsLssUUSUS2UssLsLscscscscLcsLLLLLscscsLscLLL77r7r775,                              
                .w7r7r7r7r7r77777r7r7r77777r7r777r7r7r77c;rrrrr;c777777r777777777r777r7r7r7r7r7r7r7r7r7w,                              
                .277r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7rL7       rsr7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r772.                              
                .17r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7Lc       7L7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r71,                              
                .U77r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7rsv       7s77r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r772,                              
                .17r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7LL       vL7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7rvw,                              
                .U77r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r77sL       7J77r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r77U,                              
                .5c77r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7ss       vJvr7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7rcU.                              
                 Jsvr7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r77X:       ,X77r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r777sw                               
                  5sv77r7r7r7r7r7r7r7r7r7r7r7r7r7r7r7vS7         r5c7r7r7r7r7r7r7r7r7r7r7r7r7r7r7r77vs5.                               
                   sws777777r7r7r7r7r7r7r7r7r7r7777LJ5;           ;5JL77r7r7r7r7r7r7r7r7r7r7r7r7777L2J.                                
                    :LU1svL77777777r7777v77777vvLsUw7               r1UJL7c77r777777777r777777Lcs15s:                                  
                      .;7J12121Uw21U121U12wU1U12L7,                   ,rL21Uw2w21212w2w2121212wJ7;.                                    
                            ..,.,.,.,.,.,.,...                             ...,.,.,.,.,.,.,..                                          
                                                                                                                                       
                                                                                                                                       
                                                                                                                                       
                                                                                                                                       
                                                                                                                                       
                ;BBQBBBBQ7 rB.    SBJ    :B: SBBBBBBg:  7RBBBBBQH  .OBBQBBBp  :QMQBBQZgQBBBG,  7MQBBBQMc                               
                 ,....  vB, QB   7BBBr   BB UB, ..  BB ,B5 ... ;B, BB::;;::BB iBJ:::rBR::::BB .BJ  .  cB;                              
                 ;:irrr;JB: ,BX  Bg BB  aB. XQr;rrr;BB  BG;irrr:   BQ      QB :B.    Q2    gB .BU:rrr;XBL                              
                ;BGLsJLc5B:  SB:MB   BB:B5  aBLssJss77   7sssJcGB, BM      BM :B,    B5    QB .BXvsJss7c.                              
                ;BL     ;B;   BBB:   :BBB   OB          D:     ;B; BB.,,:.,BB rB;    QE    BB .BL                                      
                 1BBBBBBQB:   .Bw     LB,    DBBBBBBBB  PBBBQBBB2  :QBBBBBBM. ;B:    BU    gB  JBBBBBBBB7
	`)
	flag.Parse()

	r := util.InitGin(util.GinLoggerHandler(util.InitGinLogger()), m)
	g := r.Group("/agent")
	{
		g.GET("/filepath", file.HandleGetFile)
		g.GET("/lastfile", file.HandleGetLastFile)
		g.POST("/shell", shell.HandleShell)
		g.GET("/fullwindow", shell.HandleFullWindow)
		g.GET("/closewindow", shell.HandleCloseWin)
		g.GET("/allwindows", shell.HandleEnumWindows)
		g.POST("/health", health.HandleHealth)
		g.POST("/update", update.HandleUpdate)
		g.GET("/version", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"code": 200,
				"data": buildVersion,
			})
		})
	}
	server := &http.Server{
		Addr:    ":5100",
		Handler: r,
	}

	go func() {
		<-global.ServerSig
		util.Logger.Println("receive interrupt signal")
		util.Logger.Println("Sleep 5 seconds to wait server finish ...")
		time.Sleep(time.Second * 5)
		if err := server.Close(); err != nil {
			util.Logger.Fatal("Server Close:", err)
		}
		global.UpdateSig <- 1
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			util.Logger.Println("Server closed under request")
		} else {
			util.Logger.Fatalf("Server closed unexpect: %s", err.Error())
		}
	}
	util.Logger.Println("Server exiting")

	<-global.UpdateSig
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		util.Logger.Fatal(err)
	}
	util.Logger.Printf("Agent Path: %s\n", dir)
	_, err = os.StartProcess(filepath.Join(dir, "agent.exe"), nil, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
	if err != nil {
		util.Logger.Errorf("Error: %s", err)
	}
	util.Logger.Println("Restart ...")
}
