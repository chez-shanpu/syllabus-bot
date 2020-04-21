package cmd

import (
	"fmt"
	"github.com/chez-shanpu/syllabus-bot/pkg/db"
	"github.com/chez-shanpu/syllabus-bot/pkg/slack"
	"github.com/chez-shanpu/syllabus-bot/pkg/syllabus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func NewBoardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "board",
		Short: "collect board info and post to slack",
		RunE:  informBoardInfo,
	}

	// flags
	flags := cmd.Flags()
	flags.StringP("endpoint", "e", "", "The endpoint for post message to slack")

	// bind flag
	_ = viper.BindPFlag("board.endpoint", flags.Lookup("endpoint"))

	// bind env vars

	_ = viper.BindEnv("board.dialect", "DB_DIALECT")
	_ = viper.BindEnv("board.host", "DB_HOST")
	_ = viper.BindEnv("board.port", "DB_PORT")
	_ = viper.BindEnv("board.user", "DB_USER")
	_ = viper.BindEnv("board.dbname", "DB_NAME")
	_ = viper.BindEnv("board.password", "DB_PASSWORD")

	// required
	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

func informBoardInfo(cmd *cobra.Command, args []string) (err error) {
	log.Print("[INFO] Start informBoardInfo()")
	endpoint := viper.GetString("board.endpoint")

	dbi := db.DBInfo{
		Dialect:  viper.GetString("board.dialect"),
		Host:     viper.GetString("board.host"),
		Port:     viper.GetString("board.port"),
		User:     viper.GetString("board.user"),
		DBName:   viper.GetString("board.dbname"),
		Password: viper.GetString("board.password"),
	}
	log.Printf("[INFO] DBInfo: %v", dbi)

	err = syllabus.InitBoardDB(dbi)
	if err != nil {
		return err
	}

	bs := syllabus.GetBoardInfoList()
	for _, b := range *bs {
		r, err := syllabus.GetBoardRecord(dbi, b.CheckSum)
		if err != nil {
			return err
		} else if r.CheckSum != "" {
			continue
		}

		log.Println(b.CheckSum)

		err = syllabus.CreateBoardRecord(dbi, b)
		if err != nil {
			return err
		}
		postStr := "新着の授業連絡だよ！\n"
		postStr += fmt.Sprintf(
			"タイトル：%s\n"+
				"掲載日：%s\n"+
				"科目名：%s\n"+
				"内容：%s\n", b.Title, b.Date, b.ClassName, b.Content)
		err = slack.PostMessage(slack.Data{Text: postStr}, endpoint)
		if err != nil {
			return err
		}
	}

	log.Print("[INFO] Successfully finished informBoardInfo()")
	return nil
}
