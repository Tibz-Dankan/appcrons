package data

import (
	"fmt"
	"math/rand"
	"time"
)

type GenTestData struct {
	usedNames     map[string]bool
	usedEmails    map[string]bool
	usedURLs      map[string]bool
	usedPasswords map[string]bool
	randGen       *rand.Rand
}

func NewGenTestData() *GenTestData {
	source := rand.NewSource(time.Now().UnixNano())
	return &GenTestData{
		usedNames:     make(map[string]bool),
		usedEmails:    make(map[string]bool),
		usedURLs:      make(map[string]bool),
		usedPasswords: make(map[string]bool),
		randGen:       rand.New(source),
	}
}

func (g *GenTestData) RandomUniqueName() string {
	firstNames := []string{"John", "Jane", "Alice", "Bob"}
	lastNames := []string{"Doe", "Smith", "Johnson", "Brown"}

	for {
		firstName := firstNames[g.randGen.Intn(len(firstNames))]
		lastName := lastNames[g.randGen.Intn(len(lastNames))]
		nanoSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
		name := fmt.Sprintf("%s%s%s", firstName, lastName, nanoSuffix)
		if !g.usedNames[name] {
			g.usedNames[name] = true
			return name
		}
	}
}

func (g *GenTestData) RandomUniqueEmail() string {
	domains := []string{"example.com", "test.com", "mail.com", "gmail.com", "outlook.com", "yahoo.com"}

	for {
		name := g.RandomUniqueName()
		nanoSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
		email := fmt.Sprintf("%s.%s@%s", name, nanoSuffix, domains[g.randGen.Intn(len(domains))])
		if !g.usedEmails[email] {
			g.usedEmails[email] = true
			return email
		}
	}
}

func (g *GenTestData) RandomUniquePassword(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"

	for {
		password := make([]byte, length)
		for i := range password {
			password[i] = chars[g.randGen.Intn(len(chars))]
		}
		nanoSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
		passStr := fmt.Sprintf("%s%s", string(password), nanoSuffix)
		if !g.usedPasswords[passStr] {
			g.usedPasswords[passStr] = true
			return passStr
		}
	}
}

func (g *GenTestData) RandomUniqueAppName() string {
	prefixes := []string{"my", "super", "cool", "awesome", "best"}
	nouns := []string{"app", "service", "tool", "platform", "project"}

	for {
		prefix := prefixes[g.randGen.Intn(len(prefixes))]
		noun := nouns[g.randGen.Intn(len(nouns))]
		nanoSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
		appName := fmt.Sprintf("%s%s%s", prefix, noun, nanoSuffix)
		if !g.usedNames[appName] {
			g.usedNames[appName] = true
			return appName
		}
	}
}

func (g *GenTestData) RandomUniqueURL() string {
	for {
		appName := g.RandomUniqueAppName()
		nanoSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
		url := fmt.Sprintf("https://%s.onrender.com/active/%s", appName, nanoSuffix)
		if !g.usedURLs[url] {
			g.usedURLs[url] = true
			return url
		}
	}
}
