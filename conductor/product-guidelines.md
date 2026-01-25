# Product Guidelines

## Communication Style
- **Tone:** Professional and Concise. The bot should provide direct, technical, and clear responses with minimal fluff.
- **Tone Rationale:** This aligns with the "privacy-first" and reliable nature of a backend-focused project, ensuring users get the information they need quickly and accurately.

## Visual & User Interface Design
- **Expressive Embeds:** Use Discord's rich embeds to provide structured information about polls, bets, and statistics.
- **State-Based Color Coding:** Utilize consistent color schemes to indicate the status of polls (e.g., green for active, red for closed, gold for resolved).
- **Iconography:** Use icons and emojis strategically to improve readability and engagement within the UI.
- **Consistency:** Maintain a professional aesthetic that is visually appealing while remaining functional and informative.

## Technical & Operational Standards
- **Privacy by Default:** Minimize data collection and ensure user data is stored securely. The system should only retain the minimum information necessary to function.
- **Performance & Latency:** The architecture must be optimized for fast response times, ensuring a fluid experience for placing bets and viewing results.
- **Architectural Purity (Hexagonal):** Adhere strictly to the Ports and Adapters pattern. The domain logic must remain isolated from external adapters (Discord, LibSQL) to facilitate testing and future portability.
- **Reliability:** Prioritize system stability and data integrity, especially during poll resolutions and statistics updates.
