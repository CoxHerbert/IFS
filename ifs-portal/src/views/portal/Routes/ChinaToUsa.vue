<template>
  <main class="route-page">
    <nav class="breadcrumbs" aria-label="Breadcrumb">
      <router-link :to="localePath('/')">Home</router-link><span>/</span><span>Routes</span><span>/</span><span>China to USA</span>
    </nav>

    <section class="hero">
      <div>
        <p class="eyebrow">CHINA–USA FREIGHT FORWARDING</p>
        <h1>Freight Forwarder from China to USA</h1>
        <p class="lead">Ocean, air, FCL, LCL, DDP and door-to-door shipping with a plan built around your cargo, destination and delivery deadline.</p>
        <div class="actions">
          <router-link :to="localePath('/contact')"><a-button type="primary" size="large">Get a Freight Quote</a-button></router-link>
          <a href="#shipping-options"><a-button size="large">Compare Shipping Options</a-button></a>
        </div>
      </div>
      <aside class="answer-card">
        <h2>Quick answer</h2>
        <p>Ocean freight is usually the practical choice for larger, less urgent shipments. Air freight fits urgent or high-value cargo. The right option depends on volume, chargeable weight, Incoterm and final delivery address.</p>
        <dl><div><dt>Modes</dt><dd>Ocean / Air</dd></div><div><dt>Load</dt><dd>FCL / LCL</dd></div><div><dt>Terms</dt><dd>EXW / FOB / CIF / DDP</dd></div></dl>
      </aside>
    </section>

    <section id="shipping-options" class="section">
      <p class="eyebrow">SHIPPING OPTIONS</p><h2>Choose the right way to ship from China to the USA</h2>
      <div class="cards">
        <article v-for="item in options" :key="item.title" class="card"><h3>{{ item.title }}</h3><p>{{ item.text }}</p><strong>{{ item.best }}</strong></article>
      </div>
    </section>

    <section class="split section">
      <div><p class="eyebrow">TRANSIT TIME</p><h2>How long does shipping from China to the USA take?</h2><p>Transit time varies by origin, destination, service level, sailing schedule, customs and inland delivery. The ranges below are planning estimates, not guaranteed schedules.</p></div>
      <div class="table-wrap"><table><thead><tr><th>Method</th><th>Typical planning range</th><th>Suitable for</th></tr></thead><tbody><tr><td>Express</td><td>3–7 days</td><td>Samples and urgent parcels</td></tr><tr><td>Air freight</td><td>5–12 days</td><td>Urgent commercial cargo</td></tr><tr><td>Ocean freight</td><td>18–40+ days</td><td>Regular and high-volume cargo</td></tr></tbody></table></div>
    </section>

    <section class="section"><p class="eyebrow">PORTS & ROUTES</p><h2>Popular China–USA shipping routes</h2><div class="route-grid"><div v-for="route in routes" :key="route"><span>{{ route }}</span></div></div><p class="note">Actual availability and timing depend on carrier schedules. Share your origin and ZIP code for a route-specific plan.</p></section>

    <section class="section process"><p class="eyebrow">PROCESS</p><h2>From pickup in China to delivery in the USA</h2><ol><li v-for="(step, index) in steps" :key="step"><span>{{ index + 1 }}</span>{{ step }}</li></ol></section>

    <section class="split section">
      <div><p class="eyebrow">DOCUMENTS</p><h2>Information needed for a useful quote</h2><ul><li>Origin and destination, including final US ZIP code</li><li>Cargo name, quantity, dimensions, weight and total CBM</li><li>Commercial invoice and packing list</li><li>Preferred Incoterm and cargo ready date</li><li>Special handling, battery, oversized or regulated cargo details</li></ul></div>
      <aside class="cta"><h2>Send your cargo details</h2><p>We will review the route, cargo profile and delivery requirement before proposing a workable shipping option.</p><router-link :to="localePath('/contact')"><a-button type="primary" size="large">Request a Tailored Quote</a-button></router-link></aside>
    </section>

    <section class="section faq"><p class="eyebrow">FAQ</p><h2>China to USA shipping questions</h2><details v-for="item in faqs" :key="item.question"><summary>{{ item.question }}</summary><p>{{ item.answer }}</p></details></section>
  </main>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { setStructuredData } from '@/utils/seo'
import { usePortalI18n } from '@/i18n'

const { localePath } = usePortalI18n()

const options = [
  { title: 'Ocean Freight', text: 'Cost-efficient transport for pallets, machinery, furniture and regular replenishment.', best: 'Best for: larger, flexible shipments' },
  { title: 'FCL Shipping', text: 'A dedicated 20GP, 40GP or 40HQ container for one shipper’s cargo.', best: 'Best for: volume and control' },
  { title: 'LCL Shipping', text: 'Share container space when cargo does not require a full container.', best: 'Best for: smaller commercial loads' },
  { title: 'Air Freight', text: 'Faster airport-to-airport or door-to-door movement for time-sensitive cargo.', best: 'Best for: urgent or high-value goods' },
  { title: 'DDP Shipping', text: 'A coordinated service covering transport, customs and agreed final delivery scope.', best: 'Best for: simplified landed delivery' },
  { title: 'Door-to-Door', text: 'Pickup in China and delivery to the consignee or fulfillment destination.', best: 'Best for: fewer handoffs' },
]
const routes = ['Shanghai → Los Angeles / Long Beach', 'Ningbo → Los Angeles / Long Beach', 'Shenzhen → Los Angeles', 'Shanghai → New York / New Jersey', 'Shenzhen → New York', 'China → Houston / Savannah']
const steps = ['Confirm cargo details and Incoterm', 'Compare mode, route, timing and cost factors', 'Pickup, consolidation or container loading', 'Export declaration and international transport', 'US customs coordination and inland delivery']
const faqs = [
  { question: 'What is the cheapest way to ship from China to the USA?', answer: 'For larger, non-urgent cargo, ocean freight is commonly more economical than air. LCL may suit smaller volumes while FCL becomes more practical as volume grows. A quote must account for origin, destination, cargo and local charges.' },
  { question: 'Should I choose FCL or LCL?', answer: 'Choose based on total CBM, cargo sensitivity, schedule and the complete landed cost—not volume alone. FCL offers dedicated container space; LCL avoids paying for unused container capacity.' },
  { question: 'Can you ship door to door or DDP?', answer: 'Door-to-door and DDP options can be assessed after reviewing the cargo, importer arrangement, destination and compliance requirements. The exact included scope should always be confirmed in the quotation.' },
  { question: 'What affects China to USA freight cost?', answer: 'Key factors include mode, weight and volume, route, season, fuel and carrier charges, Incoterm, customs requirements and pickup or final-mile distance.' },
]

onMounted(() => {
  const pageUrl = new URL(window.location.pathname, window.location.origin).href
  setStructuredData('china-usa-service', { '@context': 'https://schema.org', '@type': 'Service', name: 'Freight Forwarding from China to USA', serviceType: 'International freight forwarding', areaServed: [{ '@type': 'Country', name: 'China' }, { '@type': 'Country', name: 'United States' }], provider: { '@type': 'Organization', name: 'IFS International Logistics', url: window.location.origin }, url: pageUrl })
  setStructuredData('china-usa-breadcrumb', { '@context': 'https://schema.org', '@type': 'BreadcrumbList', itemListElement: [{ '@type': 'ListItem', position: 1, name: 'Home', item: window.location.origin }, { '@type': 'ListItem', position: 2, name: 'China to USA Shipping', item: pageUrl }] })
  setStructuredData('china-usa-faq', { '@context': 'https://schema.org', '@type': 'FAQPage', mainEntity: faqs.map((item) => ({ '@type': 'Question', name: item.question, acceptedAnswer: { '@type': 'Answer', text: item.answer } })) })
})
</script>

<style scoped>
.route-page{width:min(1180px,calc(100% - 32px));margin:auto;padding:20px 0 64px;color:#14243a}.breadcrumbs{display:flex;gap:9px;margin:8px 0 20px;color:#64748b;font-size:14px}.hero{display:grid;grid-template-columns:1.25fr .75fr;gap:24px;padding:56px;border-radius:24px;color:#fff;background:linear-gradient(135deg,#071a33,#0f65c3)}h1{max-width:760px;margin:10px 0 16px;font-size:clamp(38px,5vw,64px);line-height:1.06}h2{margin:8px 0 16px;font-size:clamp(27px,3vw,40px)}p{line-height:1.75}.lead{max-width:700px;font-size:19px;color:#dbeafe}.eyebrow{margin:0;color:#1677ff;font-weight:800;letter-spacing:.08em}.hero .eyebrow{color:#93c5fd}.actions{display:flex;gap:12px;flex-wrap:wrap;margin-top:28px}.answer-card,.cta{padding:26px;border-radius:18px;background:#fff;color:#14243a}.answer-card h2,.cta h2{font-size:24px}.answer-card dl div{display:flex;justify-content:space-between;padding:10px 0;border-top:1px solid #e5e7eb}.answer-card dd{font-weight:700}.section{padding:64px 0;border-bottom:1px solid #e5e7eb}.cards{display:grid;grid-template-columns:repeat(3,1fr);gap:16px}.card,.route-grid div{padding:24px;border:1px solid #e1e8f0;border-radius:16px;background:#fff;box-shadow:0 12px 28px rgba(16,35,63,.06)}.card h3{font-size:21px}.card strong{color:#0f65c3}.split{display:grid;grid-template-columns:1fr 1fr;gap:48px;align-items:center}.table-wrap{overflow-x:auto}table{width:100%;border-collapse:collapse}th,td{padding:14px;text-align:left;border-bottom:1px solid #dfe6ee}th{background:#edf5ff}.route-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:14px}.note{color:#64748b}.process ol{display:grid;grid-template-columns:repeat(5,1fr);gap:12px;padding:0;list-style:none}.process li{padding:18px;border-radius:14px;background:#edf5ff;font-weight:600}.process li span{display:block;margin-bottom:12px;color:#1677ff;font-size:24px}.split ul{padding-left:22px;line-height:2}.cta{background:#071a33;color:#fff}.faq details{padding:20px 0;border-top:1px solid #dfe6ee}.faq summary{cursor:pointer;font-size:18px;font-weight:700}.faq details p{max-width:850px;color:#526277}@media(max-width:800px){.hero,.split{grid-template-columns:1fr;padding:28px}.cards,.route-grid{grid-template-columns:1fr}.process ol{grid-template-columns:1fr}.section{padding:42px 0}}
</style>
