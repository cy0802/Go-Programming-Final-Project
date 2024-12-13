package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"github.com/golang/freetype"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func getImagePrompt(events []string) string {
	prompt := "將以下每一個事件生成對應的一張圖片並將它們拼貼成一張圖片:\n"
	var builder strings.Builder
	builder.WriteString(prompt)
	for _, event := range events {
		builder.WriteString(event + "\n")
	}

	prompt = builder.String()
	return prompt
}

func getTextPrompt(events []string) string {
	prompt := "用一句話描述以下這周內容，不超過30個字:\n"
	var builder strings.Builder
	builder.WriteString(prompt)
	for _, event := range events {
		builder.WriteString(event + "\n")
	}

	prompt = builder.String()
	return prompt
}

func addTextToImage(imageBytes []byte, text string) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	newHeight := height + 30
	newBound := image.Rect(0, 0, width, newHeight)
	rgba := image.NewRGBA(newBound)
	draw.Draw(rgba, bounds, img, image.Point{}, draw.Over)
	draw.Draw(rgba, image.Rect(0, height, width, newHeight), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)

	fontBytes, err := os.ReadFile("NotoSerifTC-Black.ttf")
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	c := freetype.NewContext()
	c.SetFont(font)
	c.SetFontSize(24)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.Black)

	pt := freetype.Pt(0, newHeight-5)

	_, err = c.DrawString(text, pt)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, rgba)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func CallAI(events []string) []byte {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)

	ctx := context.Background()

	prompt := getImagePrompt(events)

	image, err := client.Images.Generate(ctx, openai.ImageGenerateParams{
		Prompt:         openai.String(prompt),
		Model:          openai.F(openai.ImageModelDallE3),
		ResponseFormat: openai.F(openai.ImageGenerateParamsResponseFormatB64JSON),
		N:              openai.Int(1),
	})
	if err != nil {
		panic(err)
	}

	imageBytes, err := base64.StdEncoding.DecodeString(image.Data[0].B64JSON)
	if err != nil {
		panic(err)
	}

	prompt = getTextPrompt(events)

	stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Seed:  openai.Int(1),
		Model: openai.F(openai.ChatModelGPT4o),
	})

	var builder strings.Builder

	for stream.Next() {
		evt := stream.Current()
		if len(evt.Choices) > 0 {
			builder.WriteString(evt.Choices[0].Delta.Content)
		}
	}

	if err := stream.Err(); err != nil {
		panic(err)
	}

	text := builder.String()

	imgWithText, err := addTextToImage(imageBytes, text)
	if err != nil {
		panic(err)
	}

	return imgWithText
}
