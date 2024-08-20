package webserver

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"text/tabwriter"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (tr *Transport) exposeText() {
	tr.engine.GET("/text", func(c *gin.Context) {
		var err error
		sensors, _ := tr.listFilteredSensors(c)
		c.Header("Content-Type", "text/plain; charset=utf-8")
		tw := tabwriter.NewWriter(c.Writer, 5, 30, 1, ' ',
			tabwriter.TabIndent|tabwriter.FilterHTML,
		)
		_, err = fmt.Fprintf(tw, "Name\tMinimum\tValue\tMaximum\tError\t%s\tDescription\t\n",
			time.Now().Format("15:04:05"))
		if err != nil {
			log.Error().Err(err).Msgf("error writing text readings header: %s", err)
		}

		for i := range sensors {
			_, err = fmt.Fprintf(tw, "%s\t%.4f\t%.4f\t%.4f\t%s\t%s\t%s\t\n",
				sensors[i].Name,
				sensors[i].Minimum,
				sensors[i].Value,
				sensors[i].Maximum,
				sensors[i].Error,
				sensors[i].UpdatedAt.Format("15:04:05"),
				sensors[i].Description+" "+sensors[i].Link,
			)
			if err != nil {
				log.Error().Err(err).Msgf("error writing text readings %v %s: %s", i, sensors[i].Name, err)
			}
		}
		err = tw.Flush()
		if err != nil {
			log.Error().Err(err).Msgf("error writing text readings: %s", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.AbortWithStatus(http.StatusOK)
	})

	tr.engine.GET("/csv", func(c *gin.Context) {
		var err error
		sensors, _ := tr.listFilteredSensors(c)
		c.Header("Content-Type", "text/csv; charset=utf-8")
		tw := csv.NewWriter(c.Writer)
		var header = []string{
			"Name", "Minimum", "Value", "Maximum",
			"Error", "UpdatedAt", "Description",
		}
		err = tw.Write(header)
		if err != nil {
			log.Error().Err(err).Msgf("error writing text readings header: %s", err)
		}

		for i := range sensors {
			err = tw.Write([]string{
				sensors[i].Name,
				fmt.Sprintf("%.4f", sensors[i].Minimum),
				fmt.Sprintf("%.4f", sensors[i].Value),
				fmt.Sprintf("%.4f", sensors[i].Maximum),
				sensors[i].Error,
				sensors[i].UpdatedAt.Format("15:04:05"),
				sensors[i].Description + " " + sensors[i].Link,
			})
			if err != nil {
				log.Error().Err(err).Msgf("error writing text readings %v %s: %s", i, sensors[i].Name, err)
			}
		}
		tw.Flush()
		err = tw.Error()
		if err != nil {
			log.Error().Err(err).Msgf("error writing text readings: %s", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.AbortWithStatus(http.StatusOK)
	})
}
