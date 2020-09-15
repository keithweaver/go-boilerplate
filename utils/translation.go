package utils

type TranslationUtils struct {
	inboundMapping map[string]string
  outboundMapping map[string]string
}

func NewInstanceOfTranslationUtils() TranslationUtils {
  // Start at version header or original
  // {
  //  "original": {
  //    "conversionFunc": "",
  //    "nextVersion": ""
  //  },
  //  "2019-02-12": {
  //    "conversionFunc": ""
  //    "nextVersion": "2019-02-13"
  //  }
  // }
	inboundMapping :=
	return TranslationUtils{}
}

func (u *TranslationUtils) InboundConversion(c *gin.Context) () {

  c.Next()
}

func (u *TranslationUtils) OutboundConversion(c *gin.Context) () {
	response, exists := c.Get("response")
	if !exists {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}

	statusCode, exists := c.Get("statusCode")
	if !exists {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}

	c.JSON(statusCode, response)
}

func (u *Transaction) Nothing() () {

}
